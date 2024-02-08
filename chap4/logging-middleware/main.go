package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type LoggingClient struct {
	log *log.Logger
}

func (c LoggingClient) RoundTrip(r *http.Request) (*http.Response, error) {
	c.log.Printf("Sending a %s request to %s over %s\n", r.Method, r.URL, r.Proto)

	resp, err := http.DefaultTransport.RoundTrip(r)
	c.log.Printf("Get back a response over %s\n", resp.Proto)

	return resp, err
}
func createHttpClientWithTimeout(d time.Duration) *http.Client {
	return &http.Client{
		Timeout: d,
	}
}

func fetchRemoteResource(client *http.Client, url string) ([]byte, error) {
	r, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s URL\n", os.Args[0])
		os.Exit(1)
	}

	myTransport := LoggingClient{log.New(os.Stdout, "", log.LstdFlags)}
	client := createHttpClientWithTimeout(15 * time.Second)
	client.Transport = &myTransport

	body, err := fetchRemoteResource(client, os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching resource: %s\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "Bytes in response: %d\n", len(body))
}
