package config

import (
	"github.com/spf13/viper"
	"strings"
)

type (
	Config struct {
		Bot BotConfig
	}

	BotConfig struct {
		Debug   bool `mapstructure:"debug"`
		Timeout int  `mapstructure:"timeout"`
		Offset  int  `mapstructure:"offset"`
		Api     string
	}
)

func Init(configPath string) (*Config, error) {
	if err := parseConfigPath(configPath); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.UnmarshalKey("bot", &cfg.Bot); err != nil {
		return nil, err
	}

	parseEnv(&cfg)

	return &cfg, nil
}

func parseConfigPath(filepath string) error {
	path := strings.Split(filepath, "/")

	viper.AddConfigPath(path[0])
	viper.SetConfigName(path[1])

	return viper.ReadInConfig()
}

func parseEnv(cfg *Config) error {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	cfg.Bot.Api = viper.GetString("BOT_API")

	return nil
}
