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

type sample struct {
	Name string `json:"name"`
	Year int    `json:"year"`
	Fast bool   `json:"fast"`
}

func TestWriteFile(t *testing.T) {
	str := "Wellcome! Enjoy your development!"
	if err := WriteFile("tempfile", []byte(str)); err != nil {
		t.Errorf("WriteFile exits with error: %s", err.Error())
	}

	f, err := ioutil.ReadFile("tempfile")
	if err != nil {
		t.Error("Can not open tempfile")
	}

	if string(f) != str {
		t.Errorf("WriteFile failed\nGot:  %s\nWant: %s", string(f), str)
	}

	os.Remove("tempfile")
}

func TestGet(t *testing.T) {
	s := sample{
		Name: "Beaver",
		Year: 2017,
		Fast: true,
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		acp := strings.Join(r.Header["Accept"], ",")
		if acp != "text/html,application/json" {
			t.Errorf("Accept header not set properly, got: %s", acp)
		}

		JSON(&s).Serve(w, http.StatusOK)
	}))
	defer ts.Close()

	out := sample{}
	h := make(http.Header)
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
		t.Errorf("JSONPod.WriteFile exits with error: %s", err.Error())
	}

	out := sample{}
	if err := JSON(&out).Open("temp.json"); err != nil {
		t.Errorf("JSONPod.Open exits with error: %s", err.Error())
	}

	if out.Name != s.Name || out.Year != s.Year || out.Fast != s.Fast {
		t.Errorf("Read/Write JSON failed\n"+
			"Got  Name: %s, Year: %d, Fast: %t\n"+
			"Want Name: %s, Year: %d, Fast: %t",
			out.Name, out.Year, out.Fast,
			s.Name, s.Year, s.Fast,
		)
	}

	os.Remove("temp.json")
}

func TestServe(t *testing.T) {
	s := sample{
		Name: "Beaver",
		Year: 2017,
		Fast: true,
	}

	w := httptest.NewRecorder()
	if err := JSON(&s).Serve(w, http.StatusTeapot); err != nil {
		t.Errorf("JSONPod.Serve exits with error: %s", err.Error())
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
	if string(out) != string(want) {
		t.Errorf("JSONPod.Serve failed\nwant: %s, out: %s",
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
	if _, err := JSON(&s).Write(&buf); err != nil {
		t.Errorf("JSONPod.Write exits with error: %s", err.Error())
	}

	want, _ := json.Marshal(&s)
	if buf.String() != string(want) {
		t.Errorf("JSONPod.Write failed\nGot:  %s\n Want: %s",
			buf.String(), string(want))
	}
}
