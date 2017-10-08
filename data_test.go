package beaver

import (
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

func TestReadWriteJSON(t *testing.T) {
	s := sample{
		Name: "Beaver",
		Year: 2017,
		Fast: true,
	}

	if err := WriteJSON("temp.json", &s); err != nil {
		t.Errorf("WriteJSON exits with error: %s", err.Error())
	}

	out := sample{}
	if err := OpenJSON("temp.json", &out); err != nil {
		t.Errorf("OpenJSON exits with error: %s", err.Error())
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

func TestServeJSON(t *testing.T) {
	s := sample{
		Name: "Beaver",
		Year: 2017,
		Fast: true,
	}

	w := httptest.NewRecorder()
	if err := ServeJSON(w, http.StatusTeapot, &s); err != nil {
		t.Errorf("ServeJSON exits with error: %s", err.Error())
	}

	res := w.Result()
	defer res.Body.Close()
	ct := res.Header.Get("Content-Type")
	if ct != "application/json; charset=utf-8" {
		t.Error("ServeJSON doesn't set Content-Type header")
	}

	if res.StatusCode != http.StatusTeapot {
		t.Error("ServeJSON doesn't set status code")
	}

	out, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error("Can not read response body")
	}

	want, _ := json.Marshal(&s)
	if string(out) != string(want) {
		t.Errorf("ServeJSON failed\nwant: %s, out: %s",
			string(want), string(out))
	}
}

func TestGetJSON(t *testing.T) {
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

		ServeJSON(w, http.StatusOK, &s)
	}))
	defer ts.Close()

	out := sample{}
	h := make(http.Header)
	h.Set("Accept", "text/html")
	GetJSON(ts.URL, h, &out)

	if out.Name != s.Name || out.Year != s.Year || out.Fast != s.Fast {
		t.Errorf("Read/Write JSON failed\n"+
			"Got  Name: %s, Year: %d, Fast: %t\n"+
			"Want Name: %s, Year: %d, Fast: %t",
			out.Name, out.Year, out.Fast,
			s.Name, s.Year, s.Fast,
		)
	}
}
