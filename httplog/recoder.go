package httplog

import "net/http"

// A recorder wraps the http.ResponseWriter interface.
// It records response's status code.
type recorder struct {
	w http.ResponseWriter
	c int
}

func (w *recorder) Header() http.Header {
	return w.w.Header()
}

func (w *recorder) Write(b []byte) (int, error) {
	return w.w.Write(b)
}

func (w *recorder) WriteHeader(code int) {
	w.w.WriteHeader(code)
	w.c = code
}
