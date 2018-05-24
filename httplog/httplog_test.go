package httplog

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
	"time"
)

// regex of a log
var reg = regexp.MustCompile(`^(\w+ )?\d{4}(-\d{2}){2}T\d{2}(:\d{2}){2}\.\d+[Z+-](\d{2}:\d{2})? [\d.]+(ms|µs|ns) \d{3}(\.\d+){3} \d{3} GET /\w* .*\n$`)

func TestFile(t *testing.T) {
	n := "temp.log"
	l := File(n)
	defer os.Remove(n)

	s := "Hello world"
	fmt.Fprint(l.out, s)

	out, _ := ioutil.ReadFile(n)
	if string(out) != s {
		t.Fatalf("File failed. Got: %s, Want: %s", out, s)
	}

	// test if logger appends data to existing file instead of erasing it
	l.File(n)
	fmt.Fprint(l.out, s)

	out, _ = ioutil.ReadFile(n)
	if string(out) != s+s {
		t.Errorf("File failed, should append data instead of creating new file. Got: %s, Want: %s", out, s+s)
	}
}

func TestOutput(t *testing.T) {
	b := bytes.Buffer{}
	l := New(nil).Output(&b)
	s := "Hello world"

	fmt.Fprint(l.out, s)
	if b.String() != s {
		t.Errorf("Output failed. Got: %s, Want: %s", b.String(), s)
	}
}

func TestPrefix(t *testing.T) {
	p := "PreText"
	l := New(nil).Prefix(p)

	s := l.prefix[0].(string)
	if s != p {
		t.Errorf("Prefix failed. Got: %s, Want: %s", s, p)
	}

	l.Prefix("")
	if len(l.prefix) != 0 {
		t.Error("Prefix failed. Empty string input should lead to empty slice")
	}
}

func TestTimeFormat(t *testing.T) {
	f := time.RFC1123Z
	l := New(nil).TimeFormat(f)

	if l.timefmt != f {
		t.Errorf("TimeFormat failed. Got: %s, Want: %s", l.timefmt, f)
	}
}

func TestRecover(t *testing.T) {
	b := bytes.Buffer{}
	l := New(&b).Recover(nil)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/path", nil)

	l.Listen(willPanic).ServeHTTP(w, r)

	result := w.Result()
	defer result.Body.Close()

	c := http.StatusInternalServerError
	if result.StatusCode != c {
		t.Error("Recover failed: Default PanicHandler not set")
	}

	if len(b.String()) == 0 {
		t.Error("Recover failed: log not generated")
	}
}

func TestListen(t *testing.T) {
	b := bytes.Buffer{}
	l := New(&b).Prefix("PreText")
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/path", nil)

	l.Listen(h).ServeHTTP(w, r)

	s := b.String()
	if !reg.MatchString(s) {
		t.Error("Listen failed. output:", s)
	}
}
