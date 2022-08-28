package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	errInvalidWord     = errors.New("word is invalid")
	errInvalidLanguage = errors.New("chosen language is invalid")
	errUnknownCommand  = errors.New("unknown command")
	errDBProblem       = errors.New("internal DB problem")
	errInternalError   = errors.New("internal server error")
)

func (b *Bot) handleError(chatID int64, err error) {

	msg := tgbotapi.NewMessage(chatID, "")

	switch err {
	case errInvalidWord:
		msg.Text = b.messages.InvalidWord
	case errInvalidLanguage:
		msg.Text = b.messages.InvalidLang
	case errUnknownCommand:
		msg.Text = b.messages.UnknownCommand
	case errDBProblem:
		msg.Text = b.messages.InternalError
	case errInternalError:
		msg.Text = b.messages.InternalError
	default:
		msg.Text = b.messages.NotResponding
	}

	b.sendMessage(msg)
}
