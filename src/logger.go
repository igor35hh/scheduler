package src

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Logger interface {
	Debug(msg string, args ...any)
	Error(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
}

var _ Logger = (*noOpLogger)(nil)

type noOpLogger struct{}

func (l noOpLogger) Debug(_ string, _ ...any) {}
func (l noOpLogger) Error(_ string, _ ...any) {}
func (l noOpLogger) Info(_ string, _ ...any)  {}
func (l noOpLogger) Warn(_ string, _ ...any)  {}

var _ Logger = (*logger)(nil)

type LogLevel int

const (
	LogLevelError LogLevel = iota
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

type logger struct {
	log   *log.Logger
	level LogLevel
}

// NewLogger returns a new Logger that logs at the given level.
func NewLogger(level LogLevel) Logger {
	l := log.New(os.Stdout, "", log.LstdFlags)
	return &logger{
		log:   l,
		level: level,
	}
}

func (l *logger) Debug(msg string, args ...any) {
	if l.level < LogLevelDebug {
		return
	}
	l.log.Printf("DEBUG: %s%s\n", msg, logFormatArgs(args...))
}

func (l *logger) Error(msg string, args ...any) {
	if l.level < LogLevelError {
		return
	}
	l.log.Printf("ERROR: %s%s\n", msg, logFormatArgs(args...))
}

func (l *logger) Info(msg string, args ...any) {
	if l.level < LogLevelInfo {
		return
	}
	l.log.Printf("INFO: %s%s\n", msg, logFormatArgs(args...))
}

func (l *logger) Warn(msg string, args ...any) {
	if l.level < LogLevelWarn {
		return
	}
	l.log.Printf("WARN: %s%s\n", msg, logFormatArgs(args...))
}

func logFormatArgs(args ...any) string {
	if len(args) == 0 {
		return ""
	}
	if len(args)%2 != 0 {
		return ", " + fmt.Sprint(args...)
	}
	var pairs []string
	for i := 0; i < len(args); i += 2 {
		pairs = append(pairs, fmt.Sprintf("%s=%v", args[i], args[i+1]))
	}
	return ", " + strings.Join(pairs, ", ")
}
