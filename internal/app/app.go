package app

import (
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vadimpk/go-oxford-dictionary-sdk"
	"github.com/vadimpk/go-oxford-telegram-bot/internal/config"
	"github.com/vadimpk/go-oxford-telegram-bot/internal/repository/boltdb"
	"github.com/vadimpk/go-oxford-telegram-bot/internal/telegram"
	oxford2 "github.com/vadimpk/go-oxford-telegram-bot/pkg/oxford"
	"log"
)

func Run(configPath string) {

	cfg, err := config.Init(configPath)

	if err != nil {
		log.Fatal(err)
	}

	oxfordClient, err := oxford.NewClient(cfg.Oxford.AppID, cfg.Oxford.AppKEY)
	if err != nil {
		log.Fatal(err)
	}

	oxfordParser := oxford2.NewOxfordParser(oxfordClient)

	db, err := bolt.Open("bot.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	settingsRepository, err := boltdb.NewSettingsRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	statesRepository, err := boltdb.NewStatesRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Bot.TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = cfg.Bot.Debug

	telegramBot := telegram.NewBot(bot, oxfordParser, settingsRepository, statesRepository, &cfg.Messages)
	telegramBot.SetParseMode(cfg.Bot.ParseMode)
	if err := telegramBot.Start(cfg); err != nil {
		log.Fatal(err)
	}

}
