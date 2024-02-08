package pkgregister_data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type pkgData struct {
	Name     string    `json:"name"`
	Version  string    `json:"version"`
	Filename string    `json:"filename"`
	Bytes    io.Reader `json:"bytes"`
}

type pkgRegisterResult struct {
	Id       string `json:"id"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

func createMultiPartMessage(data pkgData) ([]byte, string, error) {
	var b bytes.Buffer
	var err error
	var fw io.Writer
	mw := multipart.NewWriter(&b)
	fw, err = mw.CreateFormField("name")
	if err != nil {
		return nil, "", err
	}
	fmt.Fprintf(fw, data.Name)

	fw, err = mw.CreateFormField("version")
	if err != nil {
		return nil, "", err
	}
	fmt.Fprintf(fw, data.Version)

	fw, err = mw.CreateFormFile("filedata", data.Filename)
	if err != nil {
		return nil, "", err
	}
	_, err = io.Copy(fw, data.Bytes)
	err = mw.Close()
	if err != nil {
		return nil, "", err
	}

	contentType := mw.FormDataContentType()
	return b.Bytes(), contentType, nil
}

func registerPackageData(url string, data pkgData) (pkgRegisterResult, error) {
	p := pkgRegisterResult{}
	payload, contentType, err := createMultiPartMessage(data)
	if err != nil {
		return p, err
	}
	render := bytes.NewReader(payload)
	r, err := http.Post(url, contentType, render)
	if err != nil {
		return p, err
	}
	defer r.Body.Close()
	respData, err := io.ReadAll(r.Body)
	if err != nil {
		return p, err
	}
	err = json.Unmarshal(respData, &p)
	return p, err
}
