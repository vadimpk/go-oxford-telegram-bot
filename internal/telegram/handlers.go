package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

const (
	startCommand = "start"
)

const (
	unknownCommandMessage = "I don't know this command"
	startCommandMessage   = "Hello there"
)

func (b *Bot) handleCommands(message *tgbotapi.Message) error {

	switch message.Command() {
	case startCommand:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)

	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, startCommandMessage)

	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, unknownCommandMessage)

	_, err := b.bot.Send(msg)
	return err
}
