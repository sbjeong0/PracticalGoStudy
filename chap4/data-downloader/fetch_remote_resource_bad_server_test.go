package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func startBadTestHttpServer() *httptest.Server {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(60 * time.Second)
				fmt.Fprintf(w, "Hello World")
			}))
	return ts
}

func fetchBadRemoteResource(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	return io.ReadAll(r.Body)
}
func TestFetchBadRemoteResource(t *testing.T) {
	ts := startBadTestHttpServer()
	defer ts.Close()
	data, err := fetchBadRemoteResource(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	expected := "Hello World"
	get := string(data)
	if expected != get {
		t.Errorf("Expected %s, got %s", expected, get)
	}
}
