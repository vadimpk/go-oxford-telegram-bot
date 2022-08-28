package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vadimpk/go-oxford-telegram-bot/internal/config"
	"github.com/vadimpk/go-oxford-telegram-bot/internal/repository"
	"github.com/vadimpk/go-oxford-telegram-bot/internal/service"
	"log"
)

type Bot struct {
	bot                *tgbotapi.BotAPI
	oxfordParser       *service.OxfordParser
	settingsRepository repository.SettingsRepository
	statesRepository   repository.StatesRepository
}

func NewBot(bot *tgbotapi.BotAPI, oxfordParser *service.OxfordParser, settingsRep repository.SettingsRepository, statesRep repository.StatesRepository) *Bot {
	return &Bot{bot: bot, oxfordParser: oxfordParser, settingsRepository: settingsRep, statesRepository: statesRep}
}

func (b *Bot) Start(cfg *config.Config) error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel(cfg)
	if err != nil {
		return err
	}

	b.handleUpdates(updates)
	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil { // If we got a message
			// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.IsCommand() {
				b.handleCommands(update.Message)
				continue
			}

			b.handleMessage(update.Message)
		}
	}
}

func (b *Bot) initUpdatesChannel(cfg *config.Config) (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(cfg.Bot.Offset)
	u.Timeout = cfg.Bot.Timeout

	return b.bot.GetUpdatesChan(u)
}
