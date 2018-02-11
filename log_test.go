package beaver

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"testing"
)

const (
	regDateTime = "[0-9][0-9][0-9][0-9]/[0-9][0-9]/[0-9][0-9] [0-9][0-9]:[0-9][0-9]:[0-9][0-9] "
	regTag      = "[A-Z ]{5}: "
	message     = "Logging message"
)

var reg, _ = regexp.Compile("^" + regDateTime + regTag + message + "\n$")

func TestOutput(t *testing.T) {
	w := new(bytes.Buffer)

	// test default Logger
	LogOutput(w)
	Debug(message)

	if !reg.Match(w.Bytes()) {
		t.Errorf("Standard output error at default Logger\nGot: %v", w)
	}

	w.Reset()

	// test new Logger
	l := NewLogger().Output(w)
	l.Info(message)

	if !reg.Match(w.Bytes()) {
		t.Errorf("Standard output error, got: %v", w)
	}
}

func TestLevelNone(t *testing.T) {
	w := new(bytes.Buffer)

	// test default Logger
	LogOutput(w)
	LogLevel(0)
	Error(message)
	Warn(message)
	Info(message)
	Debug(message)
	if w.Len() != 0 {
		t.Errorf("Level 0 default Logger failed. Output: %v", w)
	}

	// test new Logger
	l := NewLogger().Output(w).Level(0)

	l.Error(message)
	l.Warn(message)
	l.Info(message)
	l.Debug(message)
	if w.Len() != 0 {
		t.Errorf("Level 0 new Logger failed. Output: %v", w)
	}
}

func TestFatal(t *testing.T) {
	fn := "temp_log"
	f, err := os.Create(fn)
	if err != nil {
		t.Errorf("Can not create tempfile: %v", err)
	}
	defer os.Remove(fn)
	defer f.Close()

	switch os.Getenv("BEAVER_FATAL") {
	case "default":

		// the destination of standard and error output are same
		LogOutput(f)
		LogLevel(Lfatal)
		Fatal(message)
	case "new":
		l := NewLogger().Level(Lfatal).Output(f)
		l.Fatal(message)
	default:

		// test default Logger
		cmd := exec.Command(os.Args[0], "-test.run=TestFatal")
		cmd.Env = append(os.Environ(), "BEAVER_FATAL=default")
		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); !ok || e.Success() {
			t.Errorf("TestFatal subprocess default Logger failed: %v", err)
		}

		buf, err := ioutil.ReadFile(fn)
		if err != nil {
			t.Errorf(err.Error())
		}
		if !reg.Match(buf) {
			t.Errorf("Standard Logger.Fatal output not match. Got: %s", buf)
		}
		t.Log(string(buf))

		// test new Logger
		cmd = exec.Command(os.Args[0], "-test.run=TestFatal")
		cmd.Env = append(os.Environ(), "BEAVER_FATAL=new")
		err = cmd.Run()
		if e, ok := err.(*exec.ExitError); !ok || e.Success() {
			t.Errorf("TestFatal subprocess new Logger failed: %v", err)
		}

		buf, err = ioutil.ReadFile(fn)
		if err != nil {
			t.Errorf(err.Error())
		}
		if !reg.Match(buf) {
			t.Errorf("New Logger.Fatal output not match. Got: %s", buf)
		}
	}
}

func TestError(t *testing.T) {
	w := new(bytes.Buffer)

	// test default Logger
	LogOutput(w)
	LogLevel(Lerror)
	Error(message)

	if !reg.Match(w.Bytes()) {
		t.Errorf("Default Logger.Error failed. Got: %v", w)
	}

	w.Reset()

	// test new Logger
	l := NewLogger().Output(w).Level(Lerror)
	l.Error(message)

	if !reg.Match(w.Bytes()) {
		t.Errorf("New Logger.Error failed. Got: %v", w)
	}
}

func TestWarn(t *testing.T) {
	w := new(bytes.Buffer)

	// test default Logger
	LogOutput(w)
	LogLevel(Lwarn)
	Warn(message)

	if !reg.Match(w.Bytes()) {
		t.Errorf("Default Logger.Warn failed. Got: %v", w)
	}
	w.Reset()

	// test new Logger
	l := NewLogger().Output(w).Level(Lwarn)
	l.Warn(message)

	if !reg.Match(w.Bytes()) {
		t.Errorf("New Logger.Warn failed. Got: %v", w)
	}
}

func TestInfo(t *testing.T) {
	w := new(bytes.Buffer)

	// test default Logger
	LogOutput(w)
	LogLevel(Linfo)
	Info(message)

	if !reg.Match(w.Bytes()) {
		t.Errorf("Default Logger.Info failed. Got: %v", w)
	}

	w.Reset()

	// test new Logger
	l := NewLogger().Output(w).Level(Linfo)
	l.Info(message)

	if !reg.Match(w.Bytes()) {
		t.Errorf("New Logger.Info failed. Got: %v", w)
	}
}

func TestDebug(t *testing.T) {
	w := new(bytes.Buffer)

	// test default Logger
	LogOutput(w)
	LogLevel(Ldebug)
	Debug(message)

	if !reg.Match(w.Bytes()) {
		t.Errorf("Default Logger.Error failed. Got: %v", w)
	}

	w.Reset()

	// test new Logger
	l := NewLogger().Output(w).Level(Ldebug)
	l.Debug(message)

	if !reg.Match(w.Bytes()) {
		t.Errorf("New Logger.Error failed. Got: %v", w)
	}
}

func TestTags(t *testing.T) {
	w := new(bytes.Buffer)
	r, _ := regexp.Compile("^" + regDateTime + message + "\n$")

	// test empty tag
	LogOutput(w)
	LogLevel(Lall)
	LogTags(LTag{})
	defer LogTags(DefaultLogTag)

	Error(message)
	if !r.Match(w.Bytes()) {
		t.Errorf("Default Logger.Tags failed. Got: %v", w)
	}

	w.Reset()

	// test custom tag
	tag := LTag{
		Fatal: "[fatal]",
		Error: "[error]",
		Warn:  "[warn ]",
		Info:  "[info ]",
		Debug: "[debug]",
	}
	r, _ = regexp.Compile("^" + regDateTime + `\[[a-z ]{5}\] ` + message + "\n$")
	l := NewLogger().Output(w).Level(Lwarn).Tags(tag)

	l.Warn(message)
	if !r.Match(w.Bytes()) {
		t.Errorf("Default Logger.Tags failed. Got: %v", w)
	}
}
