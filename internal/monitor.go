package internal

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Monitor struct {
    fetcher     *Fetcher
    hasher      *Hasher
    store       *Store
    verbose     bool
    contentType string
	baseDir     string
}

func NewMonitor(fetcher *Fetcher, hasher *Hasher, store *Store, verbose bool, contentType string, baseDir string) *Monitor {
    if contentType != "html" {
        contentType = "js"
    }
    if baseDir == "" {
        baseDir = "data"
    }
    return &Monitor{
        fetcher:     fetcher,
        hasher:      hasher,
        store:       store,
        verbose:     verbose,
        contentType: contentType,
        baseDir:     baseDir,
    }
}

func (m *Monitor) CheckURL(rawURL string) {
	if m.verbose {
		fmt.Printf("[inf] checking %s\n", rawURL)
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		fmt.Printf("[wrn] invalid URL: %s (%v)\n", rawURL, err)
		return
	}
	canonical := parsed.Scheme + "://" + parsed.Host + parsed.Path

	content, err := m.fetcher.Fetch(rawURL)
	if err != nil {
		fmt.Printf("[err] error fetching %s: %v\n", rawURL, err)
		return
	}
	if content == "" {
		fmt.Printf("[!] Empty response for %s\n", rawURL)
		return
	}

	rawHash := m.hasher.Hash(content)
	oldHash, exists := m.store.GetLatestHash(canonical)

	if exists && oldHash == rawHash {
		if m.verbose {
			fmt.Printf("[UNCHANGED] %s (raw hash match)\n", canonical)
		}
		return
	}

	var beautified string
	switch m.contentType {
	case "html":
		beautified, err = BeautifyHTML(content)
	default:
		beautified, err = BeautifyJS(content)
	}

	if err != nil {
		if m.verbose {
			fmt.Printf("[wrn] beautify failed for %s: %v — saving raw\n", rawURL, err)
		}
		beautified = content
	}

	newHash := m.hasher.Hash(beautified)

	savedPath, saveErr := m.saveContentFile(rawURL, beautified, newHash)
	if saveErr != nil {
		fmt.Printf("[wrn] failed to save content for %s: %v\n", rawURL, saveErr)
		return
	}

	timestamp := ""
	parts := strings.Split(savedPath, string(os.PathSeparator))
	if len(parts) >= 3 {
		timestamp = parts[2]
	}

	switch {
	case !exists:
		fmt.Printf("[NEW] %s\n", canonical)
		m.store.Update(canonical, newHash, savedPath, timestamp)

	case exists && oldHash != newHash:
		fmt.Printf("[CHANGED] %s\n", canonical)
		prev := m.store.GetLastFile(canonical)

		if prev != "" {
			fmt.Printf("    Old file: %s\n", prev)
			fmt.Printf("    New file: %s\n\n", savedPath)
		} else {
			fmt.Printf("    Old file: (unknown)\n")
			fmt.Printf("    New file: %s\n\n", savedPath)
		}

		m.store.Update(canonical, newHash, savedPath, timestamp)

	default:
		if m.verbose {
			fmt.Printf("[UNCHANGED] %s (hash %s)\n", canonical, newHash)
		}
	}

	if err := m.store.Save(m.store.path); err != nil {
		fmt.Printf("[wrn] failed to save store: %v\n", err)
	}
}

func (m *Monitor) saveContentFile(rawURL, content, hash string) (string, error) {
    parsed, err := url.Parse(rawURL)
    if err != nil {
        return "", err
    }

    domain := parsed.Hostname()
    if domain == "" {
        domain = parsed.Host
    }

    ts := Timestamp()
    trimmed := strings.TrimPrefix(parsed.Path, "/")

    if trimmed == "" {
        if m.contentType == "html" {
            trimmed = "index.html"
        } else {
            trimmed = "index.js"
        }
    }

    dirParts := append([]string{m.baseDir, domain, ts, hash}, strings.Split(filepath.Dir(trimmed), "/")...)
    dirPath := filepath.Join(dirParts...)
    if err := os.MkdirAll(dirPath, 0755); err != nil {
        return "", err
    }

    fileName := filepath.Base(trimmed)
    if fileName == "" || fileName == "." || fileName == "/" {
        if m.contentType == "html" {
            fileName = "index.html"
        } else {
            fileName = "index.js"
        }
    }

    fullPath := filepath.Join(dirPath, fileName)

    if _, err := os.Stat(fullPath); err == nil {
        if m.verbose {
            fmt.Printf("[inf] version already saved for %s (hash %s)\n", rawURL, hash)
        }
        return fullPath, nil
    }

    if m.verbose {
        fmt.Printf("[inf] saving %s -> %s\n", rawURL, fullPath)
    }

    if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
        return "", err
    }
    return fullPath, nil
}

func ListDomainTimestamps(baseDir, domain string) ([]string, error) {
    domainDir := filepath.Join(baseDir, domain)
    entries, err := os.ReadDir(domainDir)
    if err != nil {
        return nil, err
    }

    var timestamps []string
    for _, e := range entries {
        if e.IsDir() {
            timestamps = append(timestamps, e.Name())
        }
    }
    sort.Strings(timestamps)
    return timestamps, nil
}