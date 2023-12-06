package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func startTestHttpServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	}))
	return ts
}

func TestFetchRemoteResource(t *testing.T) {
	ts := startTestHttpServer()
	defer ts.Close()
	expected := "Hello World"
	data, err := fetchRemoteResource(ts.URL)

	if err != nil {
		t.Fatal(err)
	}
	if expected != string(data) {
		t.Errorf("Expected response to be: %s, Got %s", expected, data)
	}
}
