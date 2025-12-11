package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"radarlance/internal"
)

var (
    inputFile   string
    singleURL   string
    outputFile  string
    threads     int
    verbose     bool
    contentType string
	quiet       bool
    dataDir     string
)

func Execute() {
	flag.StringVar(&inputFile, "i", "", "input file containing URLs (one per line)")
	flag.StringVar(&singleURL, "u", "", "single URL to check")
	flag.StringVar(&outputFile, "o", "hashes.json", "output file to store JS file hashes")
	flag.IntVar(&threads, "t", 10, "number of concurrent threads")
	flag.BoolVar(&verbose, "v", false, "enable verbose output")
	flag.StringVar(&contentType, "type", "js", "content type: js or html")
	flag.BoolVar(&quiet, "q", false, "quiet mode (no completion output)")
    flag.StringVar(&dataDir, "d", "data", "base directory for saved files and hashes.json")
	flag.Parse()

    if verbose && quiet {
        fmt.Println("[err] -v and -q cannot be used together")
        flag.Usage()
        os.Exit(1)
    }

	if inputFile == "" && singleURL == "" {
		fmt.Println("[err] you must provide either -i <file> or -u <url>")
		flag.Usage()
		os.Exit(1)
	}

    if outputFile == "" {
        outputFile = filepath.Join(dataDir, "hashes.json")
    } else if !filepath.IsAbs(outputFile) {
        outputFile = filepath.Join(dataDir, outputFile)
    }

    if err := os.MkdirAll(dataDir, 0755); err != nil {
        log.Fatalf("[err] could not create data dir %q: %v", dataDir, err)
    }

	start := time.Now()

	var urls []string
	if singleURL != "" {
		urls = append(urls, singleURL)
	} else {
		var err error
		urls, err = internal.ReadLines(inputFile)
		if err != nil {
			log.Fatalf("[err] failed to read input file: %v", err)
		}
	}

	store := internal.LoadStore(outputFile)
	fetcher := internal.NewFetcher(threads)
	hasher := internal.NewHasher()
    monitor := internal.NewMonitor(fetcher, hasher, store, verbose, contentType, dataDir)

	if len(urls) == 1 {
		monitor.CheckURL(urls[0])
	} else {
		var wg sync.WaitGroup
		for _, url := range urls {
			wg.Add(1)
			go func(u string) {
				defer wg.Done()
				monitor.CheckURL(u)
			}(url)
		}
		wg.Wait()
	}

    if err := store.Save(outputFile); err != nil {
        log.Printf("[wrn] could not save store: %v", err)
    }

    if !quiet {
        fmt.Printf("\n[inf] Completed in %v\n", time.Since(start))
    }
}
