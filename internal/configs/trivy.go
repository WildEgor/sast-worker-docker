package configs

import (
	"github.com/spf13/viper"
	"log/slog"
)

// TrivyConfig holds the main app configurations
type TrivyConfig struct {
	API   string `mapstructure:"api"`
	Token string `mapstructure:"token"`
}

func NewTrivyConfig() *TrivyConfig {
	cfg := &TrivyConfig{}

	if err := viper.UnmarshalKey("trivy", cfg); err != nil {
		slog.Error("trivy config parse error")
	}

	return cfg
}
