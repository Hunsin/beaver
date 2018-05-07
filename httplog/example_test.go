package httplog_test

import (
	"net/http"
	"time"

	"github.com/Hunsin/beaver/httplog"
)

var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// do your staff
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
})

func Example() {
	// create a log instance which writes data to stdout
	l := httplog.New(nil)

	// set output to a specific file, creating it if not existed
	l.File("http.log")

	// or, simply call httplog.File()
	l = httplog.File("http.log")

	// you can chain you configurations
	l.File("new.log").Prefix("my-app").TimeFormat(time.RFC1123)

	// use the logger as HTTP middleware
	http.ListenAndServe(":8000", l.Listen(handler))
}
