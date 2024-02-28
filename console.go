package slogger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"sync"
)

const (
	reset = "\033[0m"

	black        = 30
	red          = 31
	green        = 32
	yellow       = 33
	blue         = 34
	magenta      = 35
	cyan         = 36
	lightGray    = 37
	darkGray     = 90
	lightRed     = 91
	lightGreen   = 92
	lightYellow  = 93
	lightBlue    = 94
	lightMagenta = 95
	lightCyan    = 96
	white        = 97
)

type ConsoleHandler struct {
	handler slog.Handler
	w       io.Writer
	mu      *sync.Mutex
}

func NewConsoleHandler(
	w io.Writer,
	opts *slog.HandlerOptions,
) slog.Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	return &ConsoleHandler{
		handler: slog.NewTextHandler(w, opts),
		w:       w,
		mu:      &sync.Mutex{},
	}
}

func (h *ConsoleHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *ConsoleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ConsoleHandler{handler: h.handler.WithAttrs(attrs), w: h.w, mu: h.mu}
}

func (h *ConsoleHandler) WithGroup(name string) slog.Handler {
	return &ConsoleHandler{handler: h.handler.WithGroup(name), w: h.w, mu: h.mu}
}

func (h *ConsoleHandler) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = h.color(lightGray, level)
	case slog.LevelInfo:
		level = h.color(cyan, level)
	case slog.LevelWarn:
		level = h.color(lightBlue, level)
	case slog.LevelError:
		level = h.color(lightRed, level)
	}

	msg := h.color(white, r.Message)

	timestamp := r.Time.Format("[15:05:05.000]")

	out := strings.Builder{}

	if len(timestamp) > 0 {
		out.WriteString(timestamp)
		out.WriteString(" ")
	}
	if len(level) > 0 {
		out.WriteString(level)
		out.WriteString(" ")
	}
	if len(msg) > 0 {
		out.WriteString(msg)
		out.WriteString(" ")
	}

	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	for k, v := range fields {
		out.WriteString(k)
		out.WriteString("=")
		out.WriteString(fmt.Sprintf("%v", v))
		out.WriteString(" ")
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.w.Write([]byte(out.String()))
	return err
}

func (h *ConsoleHandler) color(colorCode int, msg string) string {
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), msg, reset)
}
