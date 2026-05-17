package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/alexvitayu/EngAIbot/internal/config"
	"github.com/lmittmann/tint"
	"gopkg.in/natefinch/lumberjack.v2"
)

func SetupLogger(cfg *config.AppConfig) (*slog.Logger, error) {
	fileLogger := &lumberjack.Logger{
		Filename:   "../log_file.txt", //путь к файлу
		MaxSize:    10,                // 10 MB
		MaxBackups: 3,                 // количество старых файлов
		MaxAge:     28,                // дней хранить стврые логи
		Compress:   true,              // сжимать старые логи
	}

	var level slog.Leveler
	var addSource bool
	var outputs []io.Writer

	if cfg.APPEnv == "development" {
		level = slog.LevelDebug
		addSource = true
		outputs = append(outputs, os.Stdout)
	} else {
		level = slog.LevelInfo
		addSource = false
	}

	outputs = append(outputs, fileLogger) // файл с ротацией

	handler := tint.NewHandler(io.MultiWriter(outputs...), &tint.Options{
		AddSource:  addSource,
		Level:      level,
		TimeFormat: time.RFC1123,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				// Преобразуем уровень в верхний регистр
				l := a.Value.String()
				return slog.String("level", strings.ToUpper(l))
			}
			return a
		},
		NoColor: true,
	})

	logger := slog.New(handler)
	return logger, nil
}
