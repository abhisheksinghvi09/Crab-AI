package config

import (
	"io"
	"log"
	"os"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

type Logger struct {
	level LogLevel
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
	fatal *log.Logger
}

var GlobalLogger *Logger

func InitLogger() {
	level := DEBUG
	if Env.Environment == "PRODUCTION" {
		level = INFO
	}

	GlobalLogger = &Logger{
		level: level,
		debug: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lshortfile),
		info:  log.New(os.Stdout, "[INFO] ", log.LstdFlags),
		warn:  log.New(os.Stdout, "[WARN] ", log.LstdFlags),
		error: log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile),
		fatal: log.New(os.Stderr, "[FATAL] ", log.LstdFlags|log.Lshortfile),
	}
}

func (l *Logger) Debug(v ...interface{}) {
	if l.level <= DEBUG {
		l.debug.Println(v...)
	}
}

func (l *Logger) Info(v ...interface{}) {
	if l.level <= INFO {
		l.info.Println(v...)
	}
}

func (l *Logger) Warn(v ...interface{}) {
	if l.level <= WARN {
		l.warn.Println(v...)
	}
}

func (l *Logger) Error(v ...interface{}) {
	if l.level <= ERROR {
		l.error.Println(v...)
	}
}

func (l *Logger) Fatal(v ...interface{}) {
	l.fatal.Println(v...)
	os.Exit(1)
}

func (l *Logger) SetOutput(w io.Writer) {
	l.debug.SetOutput(w)
	l.info.SetOutput(w)
	l.warn.SetOutput(w)
	l.error.SetOutput(w)
	l.fatal.SetOutput(w)
}

// Convenience functions
func Debug(v ...interface{}) {
	if GlobalLogger != nil {
		GlobalLogger.Debug(v...)
	}
}

func Info(v ...interface{}) {
	if GlobalLogger != nil {
		GlobalLogger.Info(v...)
	}
}

func Warn(v ...interface{}) {
	if GlobalLogger != nil {
		GlobalLogger.Warn(v...)
	}
}

func Error(v ...interface{}) {
	if GlobalLogger != nil {
		GlobalLogger.Error(v...)
	}
}

func Fatal(v ...interface{}) {
	if GlobalLogger != nil {
		GlobalLogger.Fatal(v...)
	}
}
