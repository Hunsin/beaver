package beaver

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

// WriteFile writes bytes in given pathname.
// It overwrites the file if it already exists.
func WriteFile(path string, body []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.Write(body); err != nil {
		return err
	}

	return f.Sync()
}

// A JSONPod represents
type JSONPod struct {
	v interface{}
}

// JSON returns a pointer to JSONPod which embeded with v.
// If v is nil or not a pointer, an error may returned when calling methods
// of JSONPod.
func JSON(v interface{}) *JSONPod {
	return &JSONPod{v}
}

// Get parses the JSON-encoded data from specified URL with given header
// and stores it in j. If the h is nil, a default http.Header is applied.
// "application/json" is appended to Accept header automatically.
func (j *JSONPod) Get(url string, h http.Header) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if h != nil {
		req.Header = h
	}
	req.Header.Add("Accept", "application/json")

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(j.v)
}

// Open parses JSON file from given path and stores the result in j.
func (j *JSONPod) Open(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(j.v)
}

// Post is a shorthand for *JSONPod.Send("POST", url, h)
func (j *JSONPod) Post(url string, h http.Header) (*http.Response, error) {
	return j.Send("POST", url, h)
}

// Put is a shorthand for *JSONPod.Send("PUT", url, h)
func (j *JSONPod) Put(url string, h http.Header) (*http.Response, error) {
	return j.Send("PUT", url, h)
}

// Send issues a HTTP request to specified url with given method and header.
// The request's body is JSON-encoded data of j.v. A http.Response is returned.
// It is the caller's responsibility to close the response's Body.
func (j *JSONPod) Send(method, url string, h http.Header) (*http.Response, error) {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(j.v)

	req, err := http.NewRequest(method, url, &body)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	if h != nil {
		req.Header = h
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	return (&http.Client{}).Do(req)
}

// Serve responses with JSON-encoded data of j.v to client by given status code.
// The Content-Type header is set to "application/json; charset=utf-8".
// Additional response headers must be set before calling Serve.
func (j *JSONPod) Serve(w http.ResponseWriter, code int) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	_, err := j.Write(w)
	return err
}

// Write marshals j.v and writes the result to w.
func (j *JSONPod) Write(w io.Writer) (int, error) {
	js, err := json.Marshal(j.v)
	if err != nil {
		return 0, err
	}

	return w.Write(js)
}

// WriteFile writes JSON-encoded data of j.v to a file by given path.
// It appends ".json" to the path if it doesn't end with the string.
func (j *JSONPod) WriteFile(path string) error {
	if path[len(path)-5:] != ".json" {
		path += ".json"
	}

	b, err := json.MarshalIndent(j.v, "", "\t")
	if err != nil {
		return err
	}

	return WriteFile(path, b)
}
