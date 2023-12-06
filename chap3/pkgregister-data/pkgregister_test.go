package pkgregister_data

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func packageRegHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		d := pkgRegisterResult{}
		err := r.ParseMultipartForm(5000)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		mForm := r.MultipartForm
		f := mForm.File["filedata"][0]
		d.Id = mForm.Value["name"][0] + "-" + mForm.Value["version"][0]
		d.Filename = f.Filename
		d.Size = f.Size
		jsonData, err := json.Marshal(d)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(jsonData))
	} else {
		http.Error(w, "Invalid HTTP method specified", http.StatusMethodNotAllowed)
		return
	}
}

func startTestPackageServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(packageRegHandler))
	return ts
}

func TestRegisterPackageData(t *testing.T) {
	ts := startTestPackageServer()
	defer ts.Close()
	p := pkgData{
		Name:     "mypackage",
		Version:  "0.1",
		Filename: "mypackage-0.1-tar.gz",
		Bytes:    strings.NewReader("data"),
	}

	pResult, err := registerPackageData(ts.URL, p)
	if err != nil {
		t.Fatal(err)
	}
	if pResult.Id != "mypackage-0.1" {
		t.Errorf("Expected package id to be mypackage-0.1, Got %s", pResult.Id)
	}

	if pResult.Filename != "mypackage-0.1-tar.gz" {
		t.Errorf("Expected package filename to be mypackage-0.1-tar.gz, Got %s", pResult.Filename)
		if pResult.Size != 4 {
			t.Errorf("Expected package size to be 4, Got %d", pResult.Size)
		}
	}
}
