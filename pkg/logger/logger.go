package logger

import (
	"log/slog"
	"os"
)

type LogLevel struct {
	level slog.Level
}

func (l *LogLevel) GetENV(p []byte) error {
	if len(p) == 0 {
		return nil
	}
	return l.level.UnmarshalText(p)
}

func (l *LogLevel) SetENV() ([]byte, error) {
	return l.level.MarshalText()
}

func SetLogger(level LogLevel) {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   level.level <= slog.LevelDebug,
		Level:       level.level,
		ReplaceAttr: replaceAttr,
	})

	slog.SetDefault(slog.New(handler))
}

func replaceAttr(_ []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.TimeKey:
		a.Key = "timestamp"
	case slog.LevelKey:
		a.Key = "http.log.level"
	case slog.MessageKey:
		if a.Value.String() == "" {
			return slog.Attr{}
		}
		a.Key = "http.request.message"
	}
	return a
}
