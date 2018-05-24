package httplog

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var willPanic = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	panic("handler panic!")
})

func TestHandlePanic(t *testing.T) {
	h := chainPanicHandler(willPanic, panicHandler)
	w := httptest.NewRecorder()
	h(w, nil)

	r := w.Result()
	defer r.Body.Close()

	c := http.StatusInternalServerError
	if r.StatusCode != c {
		t.Errorf("handlePanic failed: status code not set.\nGot: %d, want: %d", r.StatusCode, c)
	}

	o, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(o) != http.StatusText(c) {
		t.Errorf("handlePanic failed: response body wrong.\nGot : %s\nWant: %s", o, http.StatusText(c))
	}
}

func TestRecoder(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "beaver")
		w.Write([]byte("Hello world"))
		w.WriteHeader(http.StatusTeapot)
	}

	w := httptest.NewRecorder()
	r := &recorder{w: w}
	h(r, nil)

	if r.c != http.StatusTeapot {
		t.Error("WriteHeader failed: recorder didn't record the status code")
	}

	if w.Header().Get("Server") != "beaver" {
		t.Error("Header failed: recorder")
	}

	if w.Body.String() != "Hello world" {
		t.Error("")
	}
}
