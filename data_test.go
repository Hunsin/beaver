package beaver

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const str = "Wellcome! Enjoy your development."

type sample struct {
	Name string `json:"name"`
	Year int    `json:"year"`
	Fast bool   `json:"fast"`
}

func trimNewline(b []byte) string {
	return strings.TrimSuffix(string(b), "\n")
}

func testRequest(t *testing.T, method string, h http.Header, v interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// check method
		if r.Method != method {
			t.Errorf("Request method not match, got: %s, want: %s", r.Method, method)
		}

		want, err := json.Marshal(v)
		if err != nil {
			t.Errorf("Can not marshal v: %v", err)
		}

		// json.Encoder.Encode appends a newline charater
		want = append(want, []byte("\n")...)

		if h == nil {
			h = make(http.Header)
		}

		// GET:          Accept: "application/json"; write body
		// Others: Content-Tyep: "application/json"; check body
		if r.Method == "GET" {
			h.Add("Accept", "application/json")
			w.Write(want)
		} else {
			h.Set("Content-Type", "application/json; charset=utf-8")

			got, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Errorf("Can not read request body: %v", err)
			}
			if string(got) != string(want) {
				t.Errorf("Request Body not match.\nGot:  %s\nWant: %s", string(got), string(want))
			}
		}

		// check header
		for key := range h {
			got := strings.Join(h[key], ",")
			want := strings.Join(r.Header[key], ",")
			if got != want {
				t.Errorf("%s header not set properly.\nGot:  %s\nWant: %s", key, got, want)
			}
		}
	}
}

func TestWriteFile(t *testing.T) {
	buf := bytes.NewBuffer([]byte(str))
	if _, err := WriteFile("tempfile", buf); err != nil {
		t.Errorf("WriteFile exits with error: %v", err)
	}

	f, err := ioutil.ReadFile("tempfile")
	if err != nil {
		t.Error("Can not open tempfile")
	}

	if string(f) != str {
		t.Errorf("WriteFile failed\nGot:  %s\nWant: %s", string(f), str)
	}
}

func TestDownload(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(str))
	})
	ts := httptest.NewServer(h)

	_, err := Download(nil, ts.URL, "tempfile")
	if err != nil {
		t.Errorf("Download failed with error: %v", err)
	}
	defer os.Remove("tempfile")

	out, err := ioutil.ReadFile("tempfile")
	if err != nil {
		t.Errorf("Can not open tempfile: %v", err)
	}

	if string(out) != str {
		t.Errorf("Download failed\nGot:  %s\nWant: %s", string(out), str)
	}
}

func TestGet(t *testing.T) {
	s := sample{
		Name: "Beaver",
		Year: 2017,
		Fast: true,
	}

	hwant := make(http.Header)
	hwant.Set("Accept", "text/html")

	ts := httptest.NewServer(testRequest(t, "GET", hwant, &s))
	defer ts.Close()

	h := make(http.Header)
	out := sample{}
	h.Set("Accept", "text/html")
	JSON(&out).Get(ts.URL, h)

	if out.Name != s.Name || out.Year != s.Year || out.Fast != s.Fast {
		t.Errorf("Read/Write JSON failed\n"+
			"Got  Name: %s, Year: %d, Fast: %t\n"+
			"Want Name: %s, Year: %d, Fast: %t",
			out.Name, out.Year, out.Fast,
			s.Name, s.Year, s.Fast,
		)
	}
}

func TestOpenWriteFile(t *testing.T) {
	s := sample{
		Name: "Beaver",
		Year: 2017,
		Fast: true,
	}

	if err := JSON(&s).WriteFile("temp.json"); err != nil {
		t.Errorf("JSONPod.WriteFile exits with error: %v", err)
	}
	defer os.Remove("temp.json")

	out := sample{}
	if err := JSON(&out).Open("temp.json"); err != nil {
		t.Errorf("JSONPod.Open exits with error: %v", err)
	}

	if out.Name != s.Name || out.Year != s.Year || out.Fast != s.Fast {
		t.Errorf("Read/Write JSON failed\n"+
			"Got  Name: %s, Year: %d, Fast: %t\n"+
			"Want Name: %s, Year: %d, Fast: %t",
			out.Name, out.Year, out.Fast,
			s.Name, s.Year, s.Fast,
		)
	}
}

func TestParse(t *testing.T) {
	s := sample{
		Name: "Beaver",
		Year: 2017,
		Fast: true,
	}

	b, err := json.Marshal(&s)
	if err != nil {
		t.Errorf("Marshal struct exits with error: %v", err)
	}

	out := sample{}
	err = JSON(&out).Parse(b)
	if err != nil {
		t.Errorf("JSONPod.Parse exits with error: %v", err)
	}

	if out.Name != s.Name || out.Year != s.Year || out.Fast != s.Fast {
		t.Errorf("Read/Write JSON failed\n"+
			"Got  Name: %s, Year: %d, Fast: %t\n"+
			"Want Name: %s, Year: %d, Fast: %t",
			out.Name, out.Year, out.Fast,
			s.Name, s.Year, s.Fast,
		)
	}
}

func TestSend(t *testing.T) {
	s := sample{
		Name: "Beaver",
		Year: 2017,
		Fast: true,
	}

	hwant := make(http.Header)
	hwant.Add("Cache-Control", "no-cache")
	ts := httptest.NewServer(testRequest(t, "POST", hwant, &s))
	defer ts.Close()

	// Content-Type header will be rewrited
	h := make(http.Header)
	h.Set("Content-Type", "unintended value")
	h.Add("Cache-Control", "no-cache")
	JSON(&s).Send("POST", ts.URL, h)
}

func TestPost(t *testing.T) {
	s := sample{
		Name: "Beaver",
		Year: 2017,
		Fast: true,
	}

	ts := httptest.NewServer(testRequest(t, "POST", make(http.Header), &s))
	defer ts.Close()

	JSON(&s).Post(ts.URL, nil)
}

func TestPut(t *testing.T) {
	s := sample{
		Name: "Beaver",
		Year: 2017,
		Fast: true,
	}

	ts := httptest.NewServer(testRequest(t, "PUT", make(http.Header), &s))
	defer ts.Close()

	JSON(&s).Put(ts.URL, nil)
}

func TestServe(t *testing.T) {
	s := sample{
		Name: "Beaver",
		Year: 2017,
		Fast: true,
	}

	w := httptest.NewRecorder()
	if err := JSON(&s).Serve(w, http.StatusTeapot); err != nil {
		t.Errorf("JSONPod.Serve exits with error: %v", err)
	}

	res := w.Result()
	defer res.Body.Close()
	ct := res.Header.Get("Content-Type")
	if ct != "application/json; charset=utf-8" {
		t.Error("JSONPod.Serve doesn't set Content-Type header")
	}

	if res.StatusCode != http.StatusTeapot {
		t.Error("JSONPod.Serve doesn't set status code")
	}

	out, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error("Can not read response body")
	}

	want, _ := json.Marshal(&s)

	if trimNewline(want) != trimNewline(out) {
		t.Errorf("JSONPod.Serve failed\nWant: %s\nGot:  %s",
			string(want), string(out))
	}
}

func TestWrite(t *testing.T) {
	s := sample{
		Name: "Beaver",
		Year: 2017,
		Fast: true,
	}

	var buf bytes.Buffer
	if err := JSON(&s).Write(&buf); err != nil {
		t.Errorf("JSONPod.Write exits with error: %v", err)
	}

	want, _ := json.Marshal(&s)
	if trimNewline(buf.Bytes()) != trimNewline(want) {
		t.Errorf("JSONPod.Write failed\nGot:  %v\n Want: %s",
			buf, string(want))
	}
}
