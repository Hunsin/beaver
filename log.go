package beaver

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// The levels define whether a Logger should generate logs in different level
const (
	Lfatal = 1 << iota
	Lerror
	Lwarn
	Linfo
	Ldebug

	Lall = Lfatal | Lerror | Lwarn | Linfo | Ldebug
)

// A LTag is a set of strings which identify the level of logs
type LTag struct {
	Fatal, Error, Warn, Info, Debug string
}

var (

	// DefaultLogTag is a LTag with default strings
	DefaultLogTag = LTag{"FATAL: ", "ERROR: ", "WARN : ", "INFO : ", "DEBUG: "}

	// BracketLogTag is a LTag with "[...]" style
	BracketLogTag = LTag{"[FATAL] ", "[ERROR] ", "[WARN ] ", "[INFO ] ", "[DEBUG] "}
)

// Logger wraps two log.Loggers with different output deestination. One is standard
// output when method Debug, Info or Warn is called; another is error output where
// Error or Fatal writes data.
type Logger struct {
	o, e  *log.Logger
	level int
	lock  sync.Mutex
	t     LTag
}

// write calls lgr.Output to print tag and v to output
func (l *Logger) write(lgr *log.Logger, tag string, v ...interface{}) {
	v = append([]interface{}{tag}, v...)
	l.lock.Lock()
	defer l.lock.Unlock()
	lgr.Output(3, fmt.Sprint(v...))
}

// Fatal calls os.Exit(1) after writes given message to error output
func (l *Logger) Fatal(v ...interface{}) {
	if l.level&Lfatal != 0 {
		l.write(l.e, l.t.Fatal, v...)
	}
	os.Exit(1)
}

// Error writes given message to error output
func (l *Logger) Error(v ...interface{}) {
	if l.level&Lerror != 0 {
		l.write(l.e, l.t.Error, v...)
	}
}

// Warn writes given message to standard output
func (l *Logger) Warn(v ...interface{}) {
	if l.level&Lwarn != 0 {
		l.write(l.o, l.t.Warn, v...)
	}
}

// Info writes given message to standard output
func (l *Logger) Info(v ...interface{}) {
	if l.level&Linfo != 0 {
		l.write(l.o, l.t.Info, v...)
	}
}

// Debug writes given message to standard output
func (l *Logger) Debug(v ...interface{}) {
	if l.level&Ldebug != 0 {
		l.write(l.o, l.t.Debug, v...)
	}
}

// Flags sets the output flags of the Logger. The flag follows the standard
// package "log"
func (l *Logger) Flags(f int) *Logger {
	l.o.SetFlags(f)
	l.e.SetFlags(f)
	return l
}

// Level sets the output level of the Logger
func (l *Logger) Level(lv int) *Logger {
	l.level = lv
	return l
}

// Output sets the destination of standard and error outputs of the Logger.
// The out must not be nil. The destination of error output will be same as
// standard output if eout is nil.
func (l *Logger) Output(out, eout io.Writer) *Logger {
	if out == nil {
		panic("A nil pointer can not be used as output")
	}
	l.o.SetOutput(out)

	if eout != nil {
		l.e.SetOutput(eout)
	} else {
		l.e.SetOutput(out)
	}

	return l
}

// Prefix sets the prefix of standard and error outputs of the Logger
func (l *Logger) Prefix(p string) *Logger {
	l.o.SetPrefix(p)
	l.e.SetPrefix(p)
	return l
}

// Tags sets the tags of the Logger. To omit the tag of logs, simply
// use l.Tags(LTag{}).
func (l *Logger) Tags(t LTag) *Logger {
	l.t = t
	return l
}

var stdLgr = NewLogger()

// Fatal calls default Logger.Fatal
func Fatal(v ...interface{}) {
	if stdLgr.level&Lfatal != 0 {
		stdLgr.write(stdLgr.e, stdLgr.t.Fatal, v...)
	}
	os.Exit(1)
}

// Error calls default Logger.Error
func Error(v ...interface{}) {
	if stdLgr.level&Lerror != 0 {
		stdLgr.write(stdLgr.e, stdLgr.t.Error, v...)
	}
}

// Warn calls default Logger.Warn
func Warn(v ...interface{}) {
	if stdLgr.level&Lwarn != 0 {
		stdLgr.write(stdLgr.o, stdLgr.t.Warn, v...)
	}
}

// Info calls default Logger.Info
func Info(v ...interface{}) {
	if stdLgr.level&Linfo != 0 {
		stdLgr.write(stdLgr.o, stdLgr.t.Info, v...)
	}
}

// Debug calls default Logger.Debug
func Debug(v ...interface{}) {
	if stdLgr.level&Ldebug != 0 {
		stdLgr.write(stdLgr.o, stdLgr.t.Debug, v...)
	}
}

// LogFlags sets the flags of default Logger
func LogFlags(f int) *Logger {
	return stdLgr.Flags(f)
}

// LogLevel sets the level of default Logger
func LogLevel(lv int) *Logger {
	return stdLgr.Level(lv)
}

// LogOutput sets the destination of standard and error outputs of default Logger
func LogOutput(out, err io.Writer) *Logger {
	return stdLgr.Output(out, err)
}

// LogPrefix sets prefix of default Logger
func LogPrefix(p string) *Logger {
	return stdLgr.Prefix(p)
}

// LogTags sets the tags of default Logger
func LogTags(t LTag) *Logger {
	return stdLgr.Tags(t)
}

// NewLogger returns a new Logger
// By default, it logs at all level and the output destination is same as default Logger.
func NewLogger() *Logger {
	return &Logger{
		level: Lall,
		o:     log.New(os.Stdout, "", log.LstdFlags),
		e:     log.New(os.Stderr, "", log.LstdFlags),
		t:     DefaultLogTag,
	}
}
