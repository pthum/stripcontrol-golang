package alog

import (
	"fmt"
	"log"
	"log/slog"
	"os"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type ilog struct {
	l *slog.Logger
}

func NewLogger(name string) Logger {
	sl := slog.New(newLogHandler()).WithGroup(name)
	// sl := slog.Default().WithGroup(slog.Group(name))
	return &ilog{
		l: sl,
	}
}
func NewLogLogger(name string) *log.Logger {
	return slog.NewLogLogger(newLogHandler(), slog.LevelDebug)
}

func newLogHandler() slog.Handler {
	return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
}

func (l *ilog) Info(msg string, args ...any) {
	fmtMsg := fmt.Sprintf(msg, args...)
	l.l.Info(fmtMsg)
}

func (l *ilog) Debug(msg string, args ...any) {
	fmtMsg := fmt.Sprintf(msg, args...)
	l.l.Debug(fmtMsg)
}

func (l *ilog) Warn(msg string, args ...any) {
	fmtMsg := fmt.Sprintf(msg, args...)
	l.l.Warn(fmtMsg)
}

func (l *ilog) Error(msg string, args ...any) {
	fmtMsg := fmt.Sprintf(msg, args...)
	l.l.Error(fmtMsg)
}
