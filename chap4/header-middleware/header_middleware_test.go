package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func startTestHttpServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			w.Header().Set(k, v[0])
		}
		fmt.Fprint(w, "I am the Request Header Echo Server echoing program")
	}))
	return ts
}

func TestAddHeaderMiddleware(t *testing.T) {
	testHeaders := map[string]string{
		"X-Client-Id": "test-client",
		"X-Auth-Hash": "random$string",
	}

	client := createClient(testHeaders)
	ts := startTestHttpServer()
	defer ts.Close()

	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatalf("Expected non-nil [AU nil-JA] error, got %s", err)
	}
	for k, v := range testHeaders {
		if resp.Header.Get(k) != testHeaders[k] {
			t.Errorf("Expected header %s to be %s, got %s", k, v, resp.Header.Get(k))
		}
	}
}
