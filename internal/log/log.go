// Package log只是SDK本身自己使用，用来调试代码使用，比如输出HTTP请求和响应信息
package log

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	*log.Logger
	level Level
}

// New 返回一个Logger 指针
func New(out io.Writer, prefix string, flag int, level Level) *Logger {
	return &Logger{
		Logger: log.New(out, prefix, flag),
		level:  level,
	}
}

var (
	std   = New(os.Stdout, DebugPrefix, log.LstdFlags, DebugLevel)
	info  = New(os.Stdout, InfoPrefix, log.LstdFlags, InfoLevel)
	warn  = New(os.Stdout, WarnPrefix, log.LstdFlags, WarnLevel)
	error = New(os.Stderr, WarnPrefix, log.LstdFlags, ErrorLevel)
)

type Level int

const (
	TraceLevel Level = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
)

const (
	InfoPrefix  = "[I] "
	DebugPrefix = "[D] "
	WarnPrefix  = "[W] "
)

func (l *Logger) Info(v ...interface{}) {
	l.output(InfoLevel, v...)
}

func (l *Logger) output(level Level, v ...interface{}) {
	if l.level <= level {
		l.Logger.Println(v...)
	}
}

func (l *Logger) Debug(v ...interface{}) {
	l.output(DebugLevel, v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.output(WarnLevel, v...)
}

func Debug(v ...interface{}) {
	std.Debug(v...)
}

func Info(v ...interface{}) {
	info.Info(v...)
}

func Warn(v ...interface{}) {
	warn.Warn(v...)
}
