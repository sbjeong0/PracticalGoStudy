package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func fetchRemoteResource(client *http.Client, url string) ([]byte, error) {
	r, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

func createHttpClientWithTimeout(d time.Duration) *http.Client {
	return &http.Client{
		Timeout: d,
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s URL\n", os.Args[0])
		os.Exit(1)
	}
	client := createHttpClientWithTimeout(15 * time.Second)
	body, err := fetchRemoteResource(client, os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching resource: %s\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "%s", body)
}
