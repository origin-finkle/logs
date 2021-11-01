package wcl

import "net/http"

type customRoundTripper struct{}

func (c customRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if token != "" {
		req.Header.Add("Authorization", token)
	}
	return http.DefaultTransport.RoundTrip(req)
}
