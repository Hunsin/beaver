package beaver

import (
	"bytes"
	"regexp"
	"testing"
)

const (
	regDateTime = "[0-9][0-9][0-9][0-9]/[0-9][0-9]/[0-9][0-9] [0-9][0-9]:[0-9][0-9]:[0-9][0-9] "
	message     = "Logging message"
)

var reg, _ = regexp.Compile("^" + regDateTime + message + "\n$")

func TestOutput(t *testing.T) {
	ow := new(bytes.Buffer)
	ew := new(bytes.Buffer)

	// test default Logger
	LogOutput(ow, ew)
	Debug(message)
	Error(message)

	if !reg.Match(ow.Bytes()) {
		t.Errorf("Standard output error at default Logger\nGot: %s", ow.String())
	}
	if !reg.Match(ew.Bytes()) {
		t.Errorf("Error output error at default Logger\nGot: %s", ow.String())
	}

	ow.Reset()
	ew.Reset()

	// test new Logger
	l := NewLogger().Output(ow, ew)
	l.Info(message)
	l.Error(message)

	if !reg.Match(ow.Bytes()) {
		t.Errorf("Standard output error, got: %s", ow.String())
	}
	if !reg.Match(ew.Bytes()) {
		t.Errorf("Error output error, got: %s", ow.String())
	}
}
