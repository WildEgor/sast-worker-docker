package configs

import (
	"github.com/spf13/viper"
	"log/slog"
)

// Configurator dummy
type Configurator struct{}

func Init() *Configurator {
	c := &Configurator{}
	c.load()

	return c
}

// load Load env data from files (default: .env, .env.local)
func (c *Configurator) load() {
	// HINT: if using .env files
	//if err := godotenv.Load(".env", ".env.local"); err != nil {
	//	slog.Error("error loading envs file", slog.Any("err", err))
	//	panic(err)
	//}

	// HINT: if using .yaml files
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		slog.Error("error loading config file", slog.Any("err", err))
		return
	}
}
