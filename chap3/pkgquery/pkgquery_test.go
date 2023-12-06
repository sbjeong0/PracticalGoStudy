package pkgquery

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func startTestPackageServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `[{"name":"pkg1","version":"1.0.0"},{"name":"pkg2","version":"2.0.0"}]`)
	}))
	return ts
}

func TestFetchPackageData(t *testing.T) {
	ts := startTestPackageServer()
	defer ts.Close()

	packages, err := fetchPackageData(ts.URL)

	if err != nil {
		t.Fatal(err)
	}

	if len(packages) != 2 {
		t.Errorf("Expected 2 packages, Got %d", len(packages))
	}
}
