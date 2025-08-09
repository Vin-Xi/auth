package logger

import (
	"io"
	"log/slog"
	"os"
	"runtime/debug"
)

type ErrorLogger struct {
	*slog.Logger
}

func (l *ErrorLogger) ErrorWithStack(msg string, err error, attrs ...slog.Attr) {
	if err == nil {
		l.Error(msg)
		return
	}

	l.Error(msg, "error", err.Error(), "stack", string(debug.Stack()))
}

func NewLogger(w io.Writer) *ErrorLogger {
	jsonHandler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	return &ErrorLogger{slog.New(jsonHandler)}
}

var Log *ErrorLogger

func Init() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stderr, file)
	Log = NewLogger(mw)
}
