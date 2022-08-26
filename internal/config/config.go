package config

import (
	"github.com/spf13/viper"
	"strings"
)

type (
	Config struct {
		Bot    BotConfig
		Oxford OxfordClientConfig
	}

	BotConfig struct {
		Debug   bool `mapstructure:"debug"`
		Timeout int  `mapstructure:"timeout"`
		Offset  int  `mapstructure:"offset"`
		APIKey  string
	}

	OxfordClientConfig struct {
		AppID  string
		AppKEY string
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

	if err := parseEnv(&cfg); err != nil {
		return nil, err
	}

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

	cfg.Bot.APIKey = viper.GetString("BOT_API")

	cfg.Oxford.AppID = viper.GetString("APP_ID")
	cfg.Oxford.AppKEY = viper.GetString("APP_KEY")

	return nil
}
