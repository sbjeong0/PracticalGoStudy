package client

import "net/http"

type AddHeaderMiddleware struct {
	headers map[string]string
}

func (c AddHeaderMiddleware) RoundTrip(r *http.Request) (*http.Response, error) {
	reqCopy := r.Clone(r.Context())
	for k, v := range c.headers {
		reqCopy.Header.Add(k, v)
	}
	return http.DefaultTransport.RoundTrip(reqCopy)
}

func createClient(headers map[string]string) *http.Client {
	return &http.Client{
		Transport: &AddHeaderMiddleware{headers},
	}
}
