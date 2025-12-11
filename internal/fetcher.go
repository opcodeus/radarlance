package internal

import (
	"io"
	"net/http"
	"time"
)

type Fetcher struct {
	client *http.Client
	sem    chan struct{}
}

func NewFetcher(concurrency int) *Fetcher {
	if concurrency < 1 {
		concurrency = 1
	}
	return &Fetcher{
		client: &http.Client{Timeout: 20 * time.Second},
		sem:    make(chan struct{}, concurrency),
	}
}

func (f *Fetcher) Fetch(url string) (string, error) {
	f.sem <- struct{}{}
	defer func() { <-f.sem }()

	resp, err := f.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", &HTTPError{StatusCode: resp.StatusCode}
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

type HTTPError struct {
	StatusCode int
}

func (e *HTTPError) Error() string {
	return http.StatusText(e.StatusCode)
}