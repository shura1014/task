package tlog

import (
	"context"
	"github.com/shura1014/logger"
)

var l *logger.Logger
var ctx = context.TODO()

const (
	DebugLevel = logger.DebugLevel

	InfoLevel  = logger.InfoLevel
	WarnLevel  = logger.WarnLevel
	ErrorLevel = logger.ErrorLevel
	TEXT       = logger.TEXT
)

func init() {
	l = logger.Default("task")
}

func Info(msg any, a ...any) {
	l.DoPrint(ctx, InfoLevel, msg, logger.GetFileNameAndLine(0), a...)
}

func Debug(msg any, a ...any) {
	l.DoPrint(ctx, DebugLevel, msg, logger.GetFileNameAndLine(0), a...)
}

func Error(msg any, a ...any) {
	l.DoPrint(ctx, ErrorLevel, msg, logger.GetFileNameAndLine(0), a...)
}

func Warn(msg any, a ...any) {
	l.DoPrint(ctx, WarnLevel, msg, logger.GetFileNameAndLine(0), a...)
}

func Text(msg any, a ...any) {
	l.DoPrint(ctx, TEXT, msg, logger.GetFileNameAndLine(0), a...)
}
