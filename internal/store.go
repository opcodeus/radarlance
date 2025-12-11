package internal

import (
	"encoding/json"
	"os"
	"sync"
)

type JSVersion struct {
	Hash      string `json:"hash"`
	Path      string `json:"path"`
	Timestamp string `json:"timestamp"`
}

type JSRecord struct {
	LatestHash string      `json:"latest_hash"`
	History    []JSVersion `json:"history"`
}

type Store struct {
	Data map[string]*JSRecord `json:"-"`
	mu   sync.Mutex
	path string
}

func LoadStore(path string) *Store {
	s := &Store{
		Data: make(map[string]*JSRecord),
		path: path,
	}

	b, err := os.ReadFile(path)
	if err == nil && len(b) > 0 {
		_ = json.Unmarshal(b, &s.Data)
	}
	return s
}

func (s *Store) Save(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if path != "" {
		s.path = path
	}

	if s.path == "" {
		return nil
	}
	data, err := json.MarshalIndent(s.Data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

func (s *Store) Update(canonicalURL, hash, filePath, timestamp string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rec, ok := s.Data[canonicalURL]
	if !ok {
		rec = &JSRecord{}
		s.Data[canonicalURL] = rec
	}
	rec.LatestHash = hash
	rec.History = append(rec.History, JSVersion{
		Hash:      hash,
		Path:      filePath,
		Timestamp: timestamp,
	})
}

func (s *Store) GetLatestHash(canonicalURL string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rec, ok := s.Data[canonicalURL]
	if !ok || rec == nil {
		return "", false
	}
	return rec.LatestHash, true
}

func (s *Store) GetLastFile(canonicalURL string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	rec, ok := s.Data[canonicalURL]
	if !ok || rec == nil || len(rec.History) == 0 {
		return ""
	}
	return rec.History[len(rec.History)-1].Path
}