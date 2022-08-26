package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vadimpk/go-oxford-dictionary-sdk"
	"github.com/vadimpk/go-oxford-telegram-bot/internal/config"
	"log"
)

type Bot struct {
	bot    *tgbotapi.BotAPI
	oxford *oxford.Client
}

func NewBot(bot *tgbotapi.BotAPI, oxfordClient *oxford.Client) *Bot {
	return &Bot{bot: bot, oxford: oxfordClient}
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
