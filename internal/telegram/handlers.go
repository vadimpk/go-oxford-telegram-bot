package telegram

import (
	"fmt"
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
		return err
		// TODO: init state for the user (new command)
	}

	err, settings := b.settingsRepository.Get(message.Chat.ID)
	if err != nil {
		return err
	}

	switch state {
	case defaultState:
		text := fmt.Sprintf("secondary language: %v\nsynonyms: %v\nsentences: %v\nexamples: %v\ntranslations: %v\n", settings.SecondaryLang, settings.Synonyms, settings.Sentences, settings.Examples, settings.Translations)
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

	err, _ := b.settingsRepository.Get(message.Chat.ID)
	if err != nil {
		// initialize new settings for this user
		// TODO: use new function to init
		settings := service.NewSettings()
		if err := b.settingsRepository.Save(message.Chat.ID, settings); err != nil {
			return err
		}
	}

	err, _ = b.statesRepository.Get(message.Chat.ID)
	if err != nil {
		// initialize default state for this user
		// TODO: use new function to init
		if err := b.statesRepository.Save(message.Chat.ID, defaultState); err != nil {
			return err
		}
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, startCommandMessage)
	_, err = b.bot.Send(msg)

	return err
}

func (b *Bot) handleSetSecondaryLangCommand(chatID int64) error {
	msg := tgbotapi.NewMessage(chatID, chooseLanguageMessage)
	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}
	return b.statesRepository.Save(chatID, chooseLanguageState)
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

	err, settings := b.settingsRepository.Get(message.Chat.ID)

	msg := tgbotapi.NewMessage(message.Chat.ID, settings.SecondaryLang)
	_, err = b.bot.Send(msg)
	return err
}
