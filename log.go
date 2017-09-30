package beaver

import (
	"io"
	"log"
	"os"
)

const (
	Lfatal = 1 << iota
	Lerror
	Lwarn
	Linfo
	Ldebug
)

type Logger struct {
	o, e  *log.Logger
	level int
}

// Fatal calls os.Exit(1) after writes given message to error output
func (l *Logger) Fatal(msg string) {
	if l.level&Lfatal != 0 {
		l.e.Output(2, msg)
	}
	os.Exit(1)
}

// Error writes given message to error output
func (l *Logger) Error(msg string) {
	if l.level&Lerror != 0 {
		l.e.Output(2, msg)
	}
}

// Warn writes given message to standard output
func (l *Logger) Warn(msg string) {
	if l.level&Lwarn != 0 {
		l.o.Output(2, msg)
	}
}

// Info writes given message to standard output
func (l *Logger) Info(msg string) {
	if l.level&Linfo != 0 {
		l.o.Output(2, msg)
	}
}

// Debug writes given message to standard output
func (l *Logger) Debug(msg string) {
	if l.level&Ldebug != 0 {
		l.o.Output(2, msg)
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

// Output sets the destination of standard and error outputs of the Logger
// The out must not be nil. The destination of error output will be set to
// standard output if err is nil.
func (l *Logger) Output(out, err io.Writer) *Logger {

	if out == nil {
		panic("A nil pointer can not be used as output")
	}
	l.o.SetOutput(out)

	if err != nil {
		l.e.SetOutput(err)
		return l
	}

	l.e.SetOutput(out)
	return l
}

// Prefix sets the prefix of standard and error outputs of the Logger
func (l *Logger) Prefix(p string) *Logger {
	l.o.SetPrefix(p)
	l.e.SetPrefix(p)
	return l
}

var defaultLogger = &Logger{
	level: Lfatal | Lerror | Lwarn | Linfo | Ldebug,
	o:     log.New(os.Stdout, "", log.LstdFlags),
	e:     log.New(os.Stderr, "", log.LstdFlags),
}

// Fatal calls default Logger.Fatal
func Fatal(msg string) {
	defaultLogger.Fatal(msg)
}

// Error calls default Logger.Error
func Error(msg string) {
	defaultLogger.Error(msg)
}

// Warn calls default Logger.Warn
func Warn(msg string) {
	defaultLogger.Warn(msg)
}

// Info calls default Logger.Info
func Info(msg string) {
	defaultLogger.Info(msg)
}

// Debug calls default Logger.Debug
func Debug(msg string) {
	defaultLogger.Debug(msg)
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

// NewLogger returns a new Logger with given format.
// By default, it logs at all level and the output destination is same as default Logger.
func NewLogger(format string) *Logger {
	return &Logger{
		level: Lfatal | Lerror | Lwarn | Linfo | Ldebug,
		o:     defaultLogger.o,
		e:     defaultLogger.e,
	}
}
