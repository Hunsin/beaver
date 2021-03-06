package httplog

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

// A Logger represents an logging object that records a series of HTTP
// requests, writing data to an io.Writer. It guarantees to serialize
// access to the Writer.
type Logger struct {
	timefmt string
	prefix  []interface{}
	out     io.Writer
	ph      PanicHandler
	mu      sync.Mutex
}

// File sets the output destination to the named file. If the file
// doesn't exist, a new file with permission mode 0644 is created.
// If any error occurred when opening the file, it panics.
func (l *Logger) File(name string) *Logger {
	f, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	return l.Output(f)
}

// Output sets the output destination for the logger.
func (l *Logger) Output(w io.Writer) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	// close original file if it's present
	if old, ok := l.out.(io.Closer); ok && old != os.Stdout {
		old.Close()
	}

	l.out = w
	return l
}

// Prefix sets the output prefix for the logger.
func (l *Logger) Prefix(prefix string) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	if prefix == "" {
		l.prefix = []interface{}{}
	} else {
		l.prefix = []interface{}{prefix}
	}
	return l
}

// TimeFormat sets the time format for the logger. The format rule
// follows the function Time.Format.
func (l *Logger) TimeFormat(layout string) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.timefmt = layout
	return l
}

// Recover sets the PanicHandler of l. If h is nil, a default handler
// is applied. In other words, the handler of l can not be removed
// once it is set.
func (l *Logger) Recover(h PanicHandler) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	if h == nil {
		h = panicHandler
	}

	l.ph = h
	return l
}

// Listen is a http middleware that records the request event and
// calls the given http.Handler.
func (l *Logger) Listen(h http.Handler) http.HandlerFunc {
	if l.ph != nil {
		h = chainPanicHandler(h, l.ph)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		rec := &recorder{w, 0}
		h.ServeHTTP(rec, r)
		dur := time.Since(now)

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		v := append(l.prefix,
			now.Format(l.timefmt),
			dur,
			ip,
			rec.c,
			r.Method,
			r.URL.Path,
			r.Referer(),
			r.UserAgent(),
		)

		l.mu.Lock()
		defer l.mu.Unlock()
		fmt.Fprintln(l.out, v...)
	}
}

// File returns a new Logger which output destination is set to
// named file.
func File(name string) *Logger {
	return New(nil).File(name)
}

// New returns an initialized Logger which output destination is w.
// If w is nil, the os.Stdout is applied. The default time format
// is set to time.RFC3339Nano.
func New(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{time.RFC3339Nano, []interface{}{}, w, nil, sync.Mutex{}}
}
