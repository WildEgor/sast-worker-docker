package configs

import (
	"github.com/spf13/viper"
	"log/slog"
)

// AppConfig holds the main app configurations
type AppConfig struct {
	Name     string `mapstructure:"name"`
	HTTPPort string `mapstructure:"http_port"`
	GRPCPort string `mapstructure:"grpc_port"`
	Mode     string `mapstructure:"mode"`

	TempPath    string `mapstructure:"temp_path"`
	ScriptsPath string `mapstructure:"scripts_path"`
}

func NewAppConfig() *AppConfig {
	cfg := &AppConfig{}

	if err := viper.UnmarshalKey("app", cfg); err != nil {
		slog.Error("app config parse error")
	}

	return cfg
}

// IsProduction Check is application running in production mode
func (ac AppConfig) IsProduction() bool {
	return ac.Mode != "develop"
}
