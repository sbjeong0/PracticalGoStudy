package cmd

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func regHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		defer r.Body.Close()
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, string(data))
	} else {
		http.Error(w, "Invalid HTTP method specified", http.StatusMethodNotAllowed)
		return
	}
}

func startTestServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(regHandler))
	return ts
}
func TestPostJson(t *testing.T) {
	ts := startTestServer()
	defer ts.Close()
	resp, err := HttpPostJson(ts.URL, "{"+
		"\"name\": \"mypackage\"}")
	if err != nil {
		t.Errorf("Expected error, Got nil")
	}
	if string(resp) != "{\"name\": \"mypackage\"}" {
		t.Errorf("Expected response to be {\"name\": \"mypackage\"}, Got %s", string(resp))
	}
	fmt.Fprintf(os.Stdout, "Response: %s\n", string(resp))
}

func TestPostJsonFile(t *testing.T) {
	ts := startTestServer()
	defer ts.Close()
	b, err := os.ReadFile("test.json")
	resp, err := HttpPostJson(ts.URL, string(b))
	if err != nil {
		t.Errorf("Expected error, Got nil")
	}
	if string(resp) != "{\n  \"name\": \"test\",\n  \"version\": \"1.0.0\"\n}" {
		t.Errorf("Expected response to be {\"name\": \"test\",\"version\": \"1.0.0\"}, Got %s", string(resp))
	}
	fmt.Fprintf(os.Stdout, "Response: %s\n", string(resp))
}
