package internal

import (
	"compress/gzip"
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

func Getarch() string {
	Arch := runtime.GOARCH
	switch Arch {
	case "amd64":
		return "64"
	case "386":
		return "32"
	default:
		panic("???")
	}
}

func HttpGet(cxt context.Context, aurl string, t *http.Transport, header http.Header) (*http.Response, *time.Timer, error) {
	ctx, cancel := context.WithCancel(cxt)
	rep, err := http.NewRequestWithContext(ctx, "GET", aurl, nil)
	timer := time.AfterFunc(5*time.Second, func() {
		cancel()
	})
	if err != nil {
		return nil, nil, fmt.Errorf("HttpGet: %w", err)
	}
	if header != nil {
		rep.Header = header
	}
	rep.Header.Set("Accept", "*/*")
	rep.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")
	rep.Header.Set("Accept-Encoding", "gzip")
	c := http.Client{
		Transport: t,
	}
	reps, err := c.Do(rep)
	if err != nil {
		return reps, nil, fmt.Errorf("HttpGet: %w", err)
	}
	if reps.Header.Get("Content-Encoding") == "gzip" {
		reps.Body, err = gzip.NewReader(reps.Body)
		if err != nil {
			return reps, nil, fmt.Errorf("HttpGet: %w", err)
		}
	}
	return reps, timer, nil
}
