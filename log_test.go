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
		t.Errorf("Standard output error at default Logger\nGot: %v", ow)
	}
	if !reg.Match(ew.Bytes()) {
		t.Errorf("Error output error at default Logger\nGot: %v", ow)
	}

	ow.Reset()
	ew.Reset()

	// test new Logger
	l := NewLogger().Output(ow, ew)
	l.Info(message)
	l.Error(message)

	if !reg.Match(ow.Bytes()) {
		t.Errorf("Standard output error, got: %v", ow)
	}
	if !reg.Match(ew.Bytes()) {
		t.Errorf("Error output error, got: %v", ow)
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
		t.Errorf("Level 0 default Logger.Error failed. Error output: %v", ew)
	}
	if ow.Len() != 0 {
		t.Errorf("Level 0 default Logger failed. Output: %v", ow)
	}

	// test new Logger
	l := NewLogger().Output(ow, ew).Level(0)

	l.Error(message)
	l.Warn(message)
	l.Info(message)
	l.Debug(message)
	if ew.Len() != 0 {
		t.Errorf("Level 0 new Logger.Error failed. Error output: %v", ew)
	}
	if ow.Len() != 0 {
		t.Errorf("Level 0 new Logger failed. Output: %v", ow)
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
		LogOutput(f, nil)
		LogLevel(Lfatal)
		Fatal(message)
	case "new":
		l := NewLogger().Level(Lfatal).Output(f, nil)
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
	ow := new(bytes.Buffer)
	ew := new(bytes.Buffer)

	// test default Logger
	LogOutput(ow, ew)
	LogLevel(Lerror)
	Error(message)

	if !reg.Match(ew.Bytes()) {
		t.Errorf("Default Logger.Error failed. Got: %v", ew)
	}
	if ow.Len() != 0 {
		t.Errorf("Default Logger.Error writes to standard output.\nGot: %v", ow)
	}

	ow.Reset()
	ew.Reset()

	// test new Logger
	l := NewLogger().Output(ow, ew).Level(Lerror)
	l.Error(message)

	if !reg.Match(ew.Bytes()) {
		t.Errorf("New Logger.Error failed. Got: %v", ew)
	}
	if ow.Len() != 0 {
		t.Errorf("New Logger.Error writes to standard output.\nGot: %v", ow)
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
		t.Errorf("Default Logger.Warn failed. Got: %v", ow)
	}
	if ew.Len() != 0 {
		t.Errorf("Default Logger.Warn writes to error output.\nGot: %v", ew)
	}

	ow.Reset()
	ew.Reset()

	// test new Logger
	l := NewLogger().Output(ow, ew).Level(Lwarn)
	l.Warn(message)

	if !reg.Match(ow.Bytes()) {
		t.Errorf("New Logger.Warn failed. Got: %v", ow)
	}
	if ew.Len() != 0 {
		t.Errorf("New Logger.Warn writes to error output.\nGot: %v", ew)
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
		t.Errorf("Default Logger.Info failed. Got: %v", ow)
	}
	if ew.Len() != 0 {
		t.Errorf("Default Logger.Info writes to error output.\nGot: %v", ew)
	}

	ow.Reset()
	ew.Reset()

	// test new Logger
	l := NewLogger().Output(ow, ew).Level(Linfo)
	l.Info(message)

	if !reg.Match(ow.Bytes()) {
		t.Errorf("New Logger.Info failed. Got: %v", ow)
	}
	if ew.Len() != 0 {
		t.Errorf("New Logger.Info writes to error output.\nGot: %v", ew)
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
		t.Errorf("Default Logger.Error failed. Got: %v", ow)
	}
	if ew.Len() != 0 {
		t.Errorf("Default Logger.Error writes to error output.\nGot: %v", ew)
	}

	ow.Reset()
	ew.Reset()

	// test new Logger
	l := NewLogger().Output(ow, ew).Level(Ldebug)
	l.Debug(message)

	if !reg.Match(ow.Bytes()) {
		t.Errorf("New Logger.Error failed. Got: %v", ow)
	}
	if ew.Len() != 0 {
		t.Errorf("New Logger.Error writes to error output.\nGot: %v", ew)
	}
}
