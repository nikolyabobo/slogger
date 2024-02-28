package slogger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type logger struct {
	log   *slog.Logger
	ctx   context.Context
	level *slog.LevelVar
}

func NewJson(ctx context.Context) Logger {
	log := &logger{
		ctx:   ctx,
		level: new(slog.LevelVar),
	}
	log.log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: log.level}))

	return log
}

func NewText(ctx context.Context) Logger {
	log := &logger{
		ctx:   ctx,
		level: new(slog.LevelVar),
	}
	log.log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: log.level}))

	return log
}

func NewConsole(ctx context.Context) Logger {
	log := &logger{
		ctx:   ctx,
		level: new(slog.LevelVar),
	}
	log.log = slog.New(NewConsoleHandler(os.Stdout, &slog.HandlerOptions{Level: log.level}))

	return log
}

func (l *logger) SetLevel(lvl string) {
	switch lvl {
	case "debug":
		l.level.Set(slog.LevelDebug)
	case "warn":
		l.level.Set(slog.LevelWarn)
	case "error":
		l.level.Set(slog.LevelError)
	default:
		l.level.Set(slog.LevelInfo)
	}
}

// Error logs at LevelError.
func (l *logger) Error(msg string, args ...any) {
	l.log.ErrorContext(l.ctx, msg, args...)
}

func (l *logger) Warn(msg string, args ...any) {
	l.log.WarnContext(l.ctx, msg, args...)
}

func (l *logger) Info(msg string, args ...any) {
	l.log.InfoContext(l.ctx, msg, args...)
}

func (l *logger) Debug(msg string, args ...any) {
	l.log.DebugContext(l.ctx, msg, args...)
}

func (l *logger) Fatal(msg string, err any) {
	txt := fmt.Errorf("%s", err)
	l.log.ErrorContext(l.ctx, msg, "error", txt)
	os.Exit(-1)
}

func (l *logger) Panic(msg string, err any) {
	txt := fmt.Errorf("%v", err)
	l.log.ErrorContext(l.ctx, msg, "panic", txt)
	os.Exit(-1)
}
