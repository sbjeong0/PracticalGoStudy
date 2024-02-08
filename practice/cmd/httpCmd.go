package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type httpConfig struct {
	url      string
	method   string
	output   bool
	body     string
	bodyFile string
	upload   string
	formData arrayFlags
}

type pkgData struct {
	Name     string    `json:"name"`
	Version  string    `json:"version"`
	Filename string    `json:"filename"`
	Bytes    io.Reader `json:"bytes"`
}

func HttpGet(url string) ([]byte, error) {
	var err error
	var resp *http.Response
	resp, err = http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func HttpPostJson(url string, body string) ([]byte, error) {
	r, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	respData, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != http.StatusOK {
		return nil, err
	}
	return respData, nil
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

type pkgRegisterResult struct {
	Id       string `json:"id"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

func HttpPostForm(url string, data pkgData) ([]byte, error) {
	p := pkgRegisterResult{}
	payload, contentType, err := createMultiPartMessage(data)
	if err != nil {
		return nil, err
	}
	render := bytes.NewReader(payload)
	r, err := http.Post(url, contentType, render)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	respData, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(respData, &p)
	return payload, nil
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return fmt.Sprint(*i)
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
func HandleHttp(w io.Writer, args []string) error {
	c := httpConfig{}

	fs := flag.NewFlagSet("http", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.BoolVar(&c.output, "output", false, "Output the response body")
	fs.StringVar(&c.method, "method", "GET", "Method to call")
	fs.StringVar(&c.body, "body", "", "Body to send in the request")
	fs.StringVar(&c.bodyFile, "body-file", "", "File to read the body from")
	fs.StringVar(&c.upload, "upload", "", "Upload a file")
	fs.Var(&c.formData, "form-data", "Form data")
	fs.Usage = func() {
		var usageString = `
http: A HTTP client.
http: <options> server`
		fmt.Fprintf(w, usageString)
		fmt.Fprintln(w)
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: ")
		fs.PrintDefaults()
	}

	err := fs.Parse(args)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if fs.NArg() != 1 && c.upload == "" {
		return ErrNoServerSpecified
	}

	c.url = fs.Arg(0)
	var resp []byte
	if c.method == http.MethodGet {
		resp, err = HttpGet(c.url)
	} else if c.method == http.MethodPost && c.body != "" {
		resp, err = HttpPostJson(c.url, c.body)
	} else if c.method == http.MethodPost && c.bodyFile != "" {
		b, err := os.ReadFile(c.bodyFile)
		if err != nil {
			return err
		}
		resp, err = HttpPostJson(c.url, string(b))
	} else if c.method == http.MethodPost && c.upload != "" {
		fmt.Println("upload: ", c.formData)
		var data pkgData
		d := &c.formData
		err := json.Unmarshal(bytes.NewBufferString(d.String()).Bytes(), &data)
		fmt.Println("data: ", data)
		if err != nil {
			return err
		}
		resp, err = HttpPostForm(c.url, data)
	}

	if err != nil {
		return err
	}
	if c.output {
		file, err := os.Create("output.out")
		if err != nil {
			return err
		}

		writer := bufio.NewWriter(file)
		_, err = writer.Write(resp)
		if err != nil {
			fmt.Println("파일에 데이터를 쓰는데 오류가 발생했습니다:", err)
			return err
		}

		err = writer.Flush()
		if err != nil {
			fmt.Println("버퍼를 비우는데 오류가 발생했습니다:", err)
			return err
		}
		defer file.Close()
	}

	fmt.Println("Send HTTP Request Success!")

	return nil
}
