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

// DefaultLogTag is a LTag with default strings
var DefaultLogTag = LTag{"FATAL: ", "ERROR: ", "WARN : ", "INFO : ", "DEBUG: "}

// Logger wraps two log.Loggers with different output deestination. One is standard
// output when method Debug, Info or Warn is called; another is error output where
// Error or Fatal writes data.
type Logger struct {
	o, e  *log.Logger
	level int
	lock  sync.Mutex
	t     LTag
}

// Fatal calls os.Exit(1) after writes given message to error output
func (l *Logger) Fatal(v ...interface{}) {
	if l.level&Lfatal != 0 {
		v = append([]interface{}{l.t.Fatal}, v...)
		l.lock.Lock()
		l.e.Output(3, fmt.Sprint(v...))
		l.lock.Unlock()
	}
	os.Exit(1)
}

// Error writes given message to error output
func (l *Logger) Error(v ...interface{}) {
	if l.level&Lerror != 0 {
		v = append([]interface{}{l.t.Error}, v...)
		l.lock.Lock()
		defer l.lock.Unlock()
		l.e.Output(3, fmt.Sprint(v...))
	}
}

// Warn writes given message to standard output
func (l *Logger) Warn(v ...interface{}) {
	if l.level&Lwarn != 0 {
		v = append([]interface{}{l.t.Warn}, v...)
		l.lock.Lock()
		defer l.lock.Unlock()
		l.o.Output(3, fmt.Sprint(v...))
	}
}

// Info writes given message to standard output
func (l *Logger) Info(v ...interface{}) {
	if l.level&Linfo != 0 {
		v = append([]interface{}{l.t.Info}, v...)
		l.lock.Lock()
		defer l.lock.Unlock()
		l.o.Output(3, fmt.Sprint(v...))
	}
}

// Debug writes given message to standard output
func (l *Logger) Debug(v ...interface{}) {
	if l.level&Ldebug != 0 {
		v = append([]interface{}{l.t.Debug}, v...)
		l.lock.Lock()
		defer l.lock.Unlock()
		l.o.Output(3, fmt.Sprint(v...))
	}
}

// Flags sets the output flags of the Logger
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

var defaultLogger = NewLogger()

// Fatal calls default Logger.Fatal
func Fatal(v ...interface{}) {
	defaultLogger.Fatal(v...)
}

// Error calls default Logger.Error
func Error(v ...interface{}) {
	defaultLogger.Error(v...)
}

// Warn calls default Logger.Warn
func Warn(v ...interface{}) {
	defaultLogger.Warn(v...)
}

// Info calls default Logger.Info
func Info(v ...interface{}) {
	defaultLogger.Info(v...)
}

// Debug calls default Logger.Debug
func Debug(v ...interface{}) {
	defaultLogger.Debug(v...)
}

// LogFlags sets the
func LogFlags(f int) *Logger {
	return defaultLogger.Flags(f)
}

// LogLevel sets the level of default Logger
func LogLevel(lv int) *Logger {
	return defaultLogger.Level(lv)
}

// LogOutput sets the destination of standard and error outputs of default Logger
func LogOutput(out, err io.Writer) *Logger {
	return defaultLogger.Output(out, err)
}

// LogPrefix sets prefix of default Logger
func LogPrefix(p string) *Logger {
	return defaultLogger.Prefix(p)
}

// LogTags sets the tags of default Logger
func LogTags(t LTag) *Logger {
	return defaultLogger.Tags(t)
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
