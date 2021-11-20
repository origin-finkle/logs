package wcl

import (
	"net/http"

	"golang.org/x/sync/semaphore"
)

type customRoundTripper struct {
	sem *semaphore.Weighted
}

func (t *customRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if token != "" {
		req.Header.Add("Authorization", token)
	}
	if err := t.sem.Acquire(req.Context(), 1); err != nil {
		return nil, err
	}
	defer t.sem.Release(1)
	return http.DefaultTransport.RoundTrip(req)
}
