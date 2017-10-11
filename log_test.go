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

func TestLevelNone(t *testing.T) {
	ow := new(bytes.Buffer)
	ew := new(bytes.Buffer)

	// test default Logger
	LogOutput(ow, ew)
	LogLevel(0)
	Error(message)
	Warn(message)
	Info(message)
	Debug(message)
	if ew.Len() != 0 {
		t.Errorf("Level 0 default Logger.Error failed. Error output: %s", ew.String())
	}
	if ow.Len() != 0 {
		t.Errorf("Level 0 default Logger failed. Output: %s", ow.String())
	}

	// test new Logger
	l := NewLogger().Output(ow, ew).Level(0)

	l.Error(message)
	l.Warn(message)
	l.Info(message)
	l.Debug(message)
	if ew.Len() != 0 {
		t.Errorf("Level 0 new Logger.Error failed. Error output: %s", ew.String())
	}
	if ow.Len() != 0 {
		t.Errorf("Level 0 new Logger failed. Output: %s", ow.String())
	}
}

func TestError(t *testing.T) {
	ow := new(bytes.Buffer)
	ew := new(bytes.Buffer)

	// test default Logger
	LogOutput(ow, ew)
	LogLevel(Lerror)
	Error(message)

	if !reg.Match(ew.Bytes()) {
		t.Errorf("Default Logger.Error failed. Got: %s", ew.String())
	}
	if ow.Len() != 0 {
		t.Errorf("Default Logger.Error writes to standard output.\nGot: %s", ow.String())
	}

	ow.Reset()
	ew.Reset()

	// test new Logger
	l := NewLogger().Output(ow, ew).Level(Lerror)
	l.Error(message)

	if !reg.Match(ew.Bytes()) {
		t.Errorf("New Logger.Error failed. Got: %s", ew.String())
	}
	if ow.Len() != 0 {
		t.Errorf("New Logger.Error writes to standard output.\nGot: %s", ow.String())
	}
}

func TestWarn(t *testing.T) {
	ow := new(bytes.Buffer)
	ew := new(bytes.Buffer)

	// test default Logger
	LogOutput(ow, ew)
	LogLevel(Lwarn)
	Warn(message)

	if !reg.Match(ow.Bytes()) {
		t.Errorf("Default Logger.Warn failed. Got: %s", ow.String())
	}
	if ew.Len() != 0 {
		t.Errorf("Default Logger.Warn writes to error output.\nGot: %s", ew.String())
	}

	ow.Reset()
	ew.Reset()

	// test new Logger
	l := NewLogger().Output(ow, ew).Level(Lwarn)
	l.Warn(message)

	if !reg.Match(ow.Bytes()) {
		t.Errorf("New Logger.Warn failed. Got: %s", ow.String())
	}
	if ew.Len() != 0 {
		t.Errorf("New Logger.Warn writes to error output.\nGot: %s", ew.String())
	}
}

func TestInfo(t *testing.T) {
	ow := new(bytes.Buffer)
	ew := new(bytes.Buffer)

	// test default Logger
	LogOutput(ow, ew)
	LogLevel(Linfo)
	Info(message)

	if !reg.Match(ow.Bytes()) {
		t.Errorf("Default Logger.Info failed. Got: %s", ow.String())
	}
	if ew.Len() != 0 {
		t.Errorf("Default Logger.Info writes to error output.\nGot: %s", ew.String())
	}

	ow.Reset()
	ew.Reset()

	// test new Logger
	l := NewLogger().Output(ow, ew).Level(Linfo)
	l.Info(message)

	if !reg.Match(ow.Bytes()) {
		t.Errorf("New Logger.Info failed. Got: %s", ow.String())
	}
	if ew.Len() != 0 {
		t.Errorf("New Logger.Info writes to error output.\nGot: %s", ew.String())
	}
}

func TestDebug(t *testing.T) {
	ow := new(bytes.Buffer)
	ew := new(bytes.Buffer)

	// test default Logger
	LogOutput(ow, ew)
	LogLevel(Ldebug)
	Debug(message)

	if !reg.Match(ow.Bytes()) {
		t.Errorf("Default Logger.Error failed. Got: %s", ow.String())
	}
	if ew.Len() != 0 {
		t.Errorf("Default Logger.Error writes to error output.\nGot: %s", ew.String())
	}

	ow.Reset()
	ew.Reset()

	// test new Logger
	l := NewLogger().Output(ow, ew).Level(Ldebug)
	l.Debug(message)

	if !reg.Match(ow.Bytes()) {
		t.Errorf("New Logger.Error failed. Got: %s", ow.String())
	}
	if ew.Len() != 0 {
		t.Errorf("New Logger.Error writes to error output.\nGot: %s", ew.String())
	}
}
