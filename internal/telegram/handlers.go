package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vadimpk/go-oxford-telegram-bot/internal/service"
)

const (
	startCommand            = "start"
	setSecondaryLangCommand = "set_secondary_language"
	toggleTranslations      = "toggle_translations"
	toggleSentences         = "toggle_sentences"
	toggleSynonyms          = "toggle_synonyms"
	toggleExamples          = "toggle_examples"
)

const (
	defaultState        = "default"
	chooseLanguageState = "chooseLanguage"
)

const (
	unknownCommandMessage = "I don't know this command"
	startCommandMessage   = "Hello there"
	chooseLanguageMessage = "Type new secondary language:"
)

func (b *Bot) handleCommands(message *tgbotapi.Message) error {

	switch message.Command() {
	case startCommand:
		return b.handleStartCommand(message)
	case setSecondaryLangCommand:
		return b.handleSetSecondaryLangCommand(message.Chat.ID)
	case toggleTranslations:
		return b.handleToggleSomethingCommand(message.Chat.ID, service.ToggleTranslations)
	case toggleSynonyms:
		return b.handleToggleSomethingCommand(message.Chat.ID, service.ToggleSynonyms)
	case toggleExamples:
		return b.handleToggleSomethingCommand(message.Chat.ID, service.ToggleExamples)
	case toggleSentences:
		return b.handleToggleSomethingCommand(message.Chat.ID, service.ToggleSentences)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {

	err, state := b.statesRepository.Get(message.Chat.ID)
	if err != nil {
		err, _ := b.initStates(message.Chat.ID)
		if err != nil {
			// DB problem
			return err
		}
	}

	err, settings := b.settingsRepository.Get(message.Chat.ID)
	if err != nil {
		err, _ := b.initSettings(message.Chat.ID)
		if err != nil {
			// DB problem
			return err
		}
	}

	switch state {
	case defaultState:
		err, text := b.oxfordParser.Parse(message.Text, settings)
		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		_, err = b.bot.Send(msg)
		return err
	case chooseLanguageState:
		settings.SetSecondaryLang(message.Text)
		err := b.settingsRepository.Save(message.Chat.ID, settings)
		if err != nil {
			return err
		}
		return b.statesRepository.Save(message.Chat.ID, defaultState)
	}

	return err
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {

	_, _, err := b.initUser(message.Chat.ID)

	msg := tgbotapi.NewMessage(message.Chat.ID, startCommandMessage)
	_, err = b.bot.Send(msg)

	return err
}

func (b *Bot) handleSetSecondaryLangCommand(chatID int64) error {
	err := b.statesRepository.Save(chatID, chooseLanguageState)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(chatID, chooseLanguageMessage)
	_, err = b.bot.Send(msg)

	return err
}

func (b *Bot) handleToggleSomethingCommand(chatID int64, toggle func(s *service.Settings)) error {
	err, settings := b.settingsRepository.Get(chatID)
	if err != nil {
		return err
	}
	toggle(settings)
	return b.settingsRepository.Save(chatID, settings)
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, unknownCommandMessage)
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) initUser(chatID int64) (*service.Settings, string, error) {

	err, settings := b.settingsRepository.Get(chatID)
	if err != nil {
		err, settings = b.initSettings(chatID)
		if err != nil {
			return nil, "", err
		}
	}

	err, state := b.statesRepository.Get(chatID)
	if err != nil {
		err, state = b.initStates(chatID)
		if err != nil {
			return nil, "", err
		}
	}
	return settings, state, nil
}

func (b *Bot) initSettings(chatID int64) (error, *service.Settings) {
	settings := service.NewSettings()
	if err := b.settingsRepository.Save(chatID, settings); err != nil {
		return err, nil
	}
	return nil, settings
}

func (b *Bot) initStates(chatID int64) (error, string) {
	if err := b.statesRepository.Save(chatID, defaultState); err != nil {
		return err, ""
	}
	return nil, defaultState
}
