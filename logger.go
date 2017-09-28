package beaver

import (
	"io"
	"os"
	"sync"
)

const (
	Lfatal = 1 << iota
	Lerror
	Lwarn
	Linfo
	Ldebug
)

type Logger struct {
	mu     sync.Mutex
	ow     io.Writer
	ew     io.Writer
	format string
	level  int
}

func (l *Logger) log(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
}

func (l *Logger) Fatal(v ...interface{}) {
	if l.level&Lfatal != 0 {
		// fmt.Fprintf(l.ew, l.format, v)
	}
	os.Exit(1)
}

func (l *Logger) Error(v ...interface{}) {
	if l.level&Lerror != 0 {
		// fmt.Fprintf(l.ew, l.format, v)
	}
}

func (l *Logger) Warn(v ...interface{}) {
	if l.level&Lwarn != 0 {
		// fmt.Fprintf(l.ow, l.format, v)
	}
}

func (l *Logger) Info(v ...interface{}) {
	if l.level&Linfo != 0 {
		// fmt.Fprintf(l.ow, l.format, v)
	}
}

func (l *Logger) Debug(v ...interface{}) {
	if l.level&Ldebug != 0 {
		// fmt.Fprintf(l.ow, l.format, v)
	}
}

func (l *Logger) Level(lv int) *Logger {
	l.level = lv
	return l
}

func (l *Logger) Output(out, err io.Writer) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	if out == nil {
		panic("A nil pointer can not be used as output")
	}
	l.ow = out

	if err != nil {
		l.ew = err
		return l
	}
	l.ew = out

	return l
}

var defaultLogger = &Logger{
	level: Lfatal | Lerror | Lwarn | Linfo | Ldebug,
	ow:    os.Stdout,
	ew:    os.Stderr,
}

func Fatal(v ...interface{}) {
	defaultLogger.Fatal(v)
}

func Error(v ...interface{}) {
	defaultLogger.Error(v)
}

func Warn(v ...interface{}) {
	defaultLogger.Warn(v)
}

func Info(v ...interface{}) {
	defaultLogger.Info(v)
}

func Debug(v ...interface{}) {
	defaultLogger.Debug(v)
}

func LogFormat(f string) {
	defaultLogger.format = f
}

func LogLevel(lv int) {
	defaultLogger.Level(lv)
}

// LogOutput sets the standard output and error output of default Logger
func LogOutput(out, err io.Writer) {
	defaultLogger.Output(out, err)
}

// NewLogger returns a new Logger with given format.
// By default, it logs at all level and the output destination is same as default Logger.
func NewLogger(format string) *Logger {
	return &Logger{
		format: format,
		level:  Lfatal | Lerror | Lwarn | Linfo | Ldebug,
		ow:     defaultLogger.ow,
		ew:     defaultLogger.ew,
	}
}
