package httplog

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

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
