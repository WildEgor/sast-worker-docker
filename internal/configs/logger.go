package configs

import (
	"github.com/spf13/viper"
	"log/slog"
	"os"
)

var logLevelToSlogLevel = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
}

var (
	LogJsonFormat string = "json"
)

// LoggerConfig holds logger configurations
type LoggerConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

func NewLoggerConfig() *LoggerConfig {
	cfg := &LoggerConfig{}

	if err := viper.UnmarshalKey("logger", cfg); err != nil {
		slog.Error("app logger parse error", slog.Any("err", err))
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevelToSlogLevel[cfg.Level],
	}))
	if cfg.Format == LogJsonFormat {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevelToSlogLevel[cfg.Level],
		}))
	}
	slog.SetDefault(logger)

	return cfg
}
