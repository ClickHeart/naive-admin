package log

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

type Logger struct {
	l     *log.Logger
	level LogLevel
}

var logFmt = log.Ldate | log.Ltime | log.LUTC | log.Lshortfile | log.Lmsgprefix

var logger = Logger{
	l:     log.New(os.Stdout, "", logFmt),
	level: LevelDebug,
}

func SetLevel(s string) {
	Info("set log Level:", s)
	level := strings.ToLower(s)
	switch level {
	case "debug":
		logger.level = LevelDebug
	case "info":
		logger.level = LevelInfo
	case "warn":
		logger.level = LevelWarn
	case "error":
		logger.level = LevelError
	}
}

func Debug(args ...any) {
	if logger.level <= LevelDebug {
		args = append([]any{"[Debug]"}, args...)
		s := color.CyanString(fmt.Sprintln(args...))
		logger.l.Output(2, s)
	}
}

func Info(args ...any) {
	if logger.level <= LevelInfo {
		args = append([]any{"[Info]"}, args...)
		s := color.GreenString(fmt.Sprintln(args...))
		logger.l.Output(2, s)
	}
}

func Warn(args ...any) {
	if logger.level <= LevelWarn {
		args = append([]any{"[Warn]"}, args...)
		s := color.MagentaString(fmt.Sprintln(args...))
		logger.l.Output(2, s)
	}
}

func Error(args ...any) {
	if logger.level <= LevelError {
		args = append([]any{"[Error]"}, args...)
		s := color.RedString(fmt.Sprintln(args...))
		logger.l.Output(2, s)
	}
}

func Panic(args ...any) {
	if logger.level <= LevelError {
		s := color.RedString(fmt.Sprintln(args...))
		sp := color.RedString("[Panic] ") + s
		logger.l.Output(2, sp)
		panic(fmt.Sprintln(args...))
	}
}

func Fatal(args ...any) {
	if logger.level <= LevelError {
		args = append([]any{"[Fatal]"}, args...)
		s := color.RedString(fmt.Sprintln(args...))
		logger.l.Output(2, s)
		os.Exit(1)
	}
}
