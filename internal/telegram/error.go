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
		msg.Text = "Cannot find information about this word"
	case errInvalidLanguage:
		msg.Text = "This language is not available to choose"
	case errUnknownCommand:
		msg.Text = "You have entered the command which I can't handle"
	case errDBProblem:
		msg.Text = "There seems to be an internal error, please, try using the bot later"
	case errInternalError:
		msg.Text = "There seems to be an internal error, please, try using the bot later"
	default:
		msg.Text = "The bot is not responding, please try again later"
	}

	b.sendMessage(msg)
}
