package httplog

import "net/http"

// A PanicHandler handles HTTP requests and panics recovered from
// http handlers.
type PanicHandler func(http.ResponseWriter, *http.Request, interface{})

// panicHandler responses with status InternalServerError to client.
func panicHandler(w http.ResponseWriter, r *http.Request, _ interface{}) {
	code := http.StatusInternalServerError
	w.WriteHeader(code)
	w.Write([]byte(http.StatusText(code)))
}

func chainPanicHandler(h http.Handler, ph PanicHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rcv := recover(); rcv != nil {
				ph(w, r, rcv)
			}
		}()

		h.ServeHTTP(w, r)
	}
}

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
