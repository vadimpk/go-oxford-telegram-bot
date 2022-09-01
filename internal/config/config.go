package config

import (
	"github.com/spf13/viper"
	"strings"
)

type (
	Config struct {
		Bot      BotConfig
		Oxford   OxfordClientConfig
		Messages Messages
	}

	BotConfig struct {
		Debug     bool   `mapstructure:"debug"`
		Timeout   int    `mapstructure:"timeout"`
		Offset    int    `mapstructure:"offset"`
		ParseMode string `mapstructure:"parse_mode"`
		TOKEN     string
	}

	OxfordClientConfig struct {
		AppID  string
		AppKEY string
	}

	Messages struct {
		Responses `mapstructure:"responses"`
		Errors    `mapstructure:"errors"`
	}

	Responses struct {
		Start             string `mapstructure:"start"`
		ChooseLang        string `mapstructure:"choose_language"`
		ChooseLangSuccess string `mapstructure:"choose_language_success"`
		SettingOn         string `mapstructure:"setting_on"`
		SettingOff        string `mapstructure:"setting_off"`
		Help              string `mapstructure:"help"`
		Settings          string `mapstructure:"settings"`
	}

	Errors struct {
		InvalidWord    string `mapstructure:"invalid_word"`
		InvalidLang    string `mapstructure:"invalid_language"`
		UnknownCommand string `mapstructure:"unknown_command"`
		InternalError  string `mapstructure:"internal_error"`
		NotResponding  string `mapstructure:"not_responding"`
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
	if err := viper.UnmarshalKey("messages", &cfg.Messages); err != nil {
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

	cfg.Bot.TOKEN = viper.GetString("BOT_API_TOKEN")
	cfg.Oxford.AppID = viper.GetString("APP_ID")
	cfg.Oxford.AppKEY = viper.GetString("APP_KEY")
	//
	//cfg.Bot.TOKEN = os.Getenv("BOT_API_TOKEN")
	//cfg.Oxford.AppID = os.Getenv("APP_ID")
	//cfg.Oxford.AppKEY = os.Getenv("APP_KEY")

	return nil
}
