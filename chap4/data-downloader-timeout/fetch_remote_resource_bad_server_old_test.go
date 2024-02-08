package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func startBadTestHttpServerV1() *httptest.Server {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(30 * time.Second)
				fmt.Fprintf(w, "Hello World")
			}))
	return ts
}
func TestFetchBadRemoteResourceV1(t *testing.T) {
	ts := startBadTestHttpServerV1()
	defer ts.Close()
	client := createHttpClientWithTimeout(200 * time.Millisecond)
	_, err := fetchRemoteResource(client, ts.URL)
	if err != nil {
		t.Fatal("Expected non-nil error")
	}

	if !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Errorf("Expected context deadline exceeded, got %s", err)
	}
}
