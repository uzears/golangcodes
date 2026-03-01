package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type Logger interface {
	Info(msg string, fields ...any)
	Error(msg string, fields ...any)
	Debug(msg string, fields ...any)
	Warn(msg string, fields ...any)
	With(fields ...any) Logger
}

type zeroLogger struct {
	log zerolog.Logger
}

func New(level string) Logger {
	// Use Unix timestamps for performance and log pipelines
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	parsedLevel := zerolog.InfoLevel
	if l, err := zerolog.ParseLevel(strings.ToLower(strings.TrimSpace(level))); err == nil {
		parsedLevel = l
	}

	z := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger().
		Level(parsedLevel)

	return &zeroLogger{log: z}
}

func (l *zeroLogger) Info(msg string, fields ...any) {
	l.log.Info().
		Fields(toMap(fields)).
		Msg(msg)
}

func (l *zeroLogger) Error(msg string, fields ...any) {
	l.log.Error().
		Fields(toMap(fields)).
		Msg(msg)
}

func (l *zeroLogger) Debug(msg string, fields ...any) {
	l.log.Debug().
		Fields(toMap(fields)).
		Msg(msg)
}

func (l *zeroLogger) Warn(msg string, fields ...any) {
	l.log.Warn().
		Fields(toMap(fields)).
		Msg(msg)
}

func (l *zeroLogger) With(fields ...any) Logger {
	return &zeroLogger{
		log: l.log.With().
			Fields(toMap(fields)).
			Logger(),
	}
}

/*
toMap converts key-value pairs into a map.
Example:

	"user_id", "123", "email", "a@b.com"
*/
func toMap(fields []any) map[string]any {
	m := make(map[string]any)
	for i := 0; i < len(fields); i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			continue
		}
		if i+1 < len(fields) {
			m[key] = fields[i+1]
		}
	}
	return m
}
