package slogger

type Logger interface {
	Error(msg string, args ...any)
	Warn(msg string, args ...any)
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
	Fatal(msg string, err any)
	Panic(msg string, err any)
	SetLevel(lvl string)
}
