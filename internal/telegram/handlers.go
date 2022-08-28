package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/vadimpk/go-oxford-telegram-bot/internal/service"
	"strings"
)

const (
	startCommand            = "start"
	helpCommand             = "help"
	settingsCommand         = "settings"
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

func langKeyboard() tgbotapi.ReplyKeyboardMarkup {
	var keyboard = tgbotapi.NewReplyKeyboard()
	count := 0
	row := 0
	for lang := range service.Languages {
		if count == 2 {
			row++
			count = 0
		}
		if count == 0 {
			keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow())
		}

		keyboard.Keyboard[row] = append(keyboard.Keyboard[row], tgbotapi.NewKeyboardButton(lang))
		count++
	}
	return keyboard
}

func (b *Bot) handleCommands(message *tgbotapi.Message) error {

	switch message.Command() {
	case startCommand:
		return b.handleStartCommand(message.Chat.ID)
	case helpCommand:
		return b.handleHelpCommand(message.Chat.ID)
	case settingsCommand:
		return b.handleSettingsCommand(message.Chat.ID)
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
		return errUnknownCommand
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {

	err, state := b.statesRepository.Get(message.Chat.ID)
	if err != nil {
		if err, _ := b.initStates(message.Chat.ID); err != nil {
			return errDBProblem
		}
	}

	err, settings := b.settingsRepository.Get(message.Chat.ID)
	if err != nil {
		if err, _ := b.initSettings(message.Chat.ID); err != nil {
			return errDBProblem
		}
	}

	switch state {
	case defaultState:
		err, text := b.oxfordParser.Parse(message.Text, settings)
		if err != nil {
			return errInternalError
		}
		if text == "" {
			return errInvalidWord
		}

		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		b.sendMessage(msg)

	case chooseLanguageState:
		if err := settings.SetSecondaryLang(message.Text); err != nil {
			return errInvalidLanguage
		}
		if err := b.settingsRepository.Save(message.Chat.ID, settings); err != nil {
			return errDBProblem
		}
		if err := b.statesRepository.Save(message.Chat.ID, defaultState); err != nil {
			return errDBProblem
		}
		msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.ChooseLangSuccess)
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		b.sendMessage(msg)
	}
	return nil
}

func (b *Bot) handleStartCommand(chatID int64) error {

	if _, _, err := b.initUser(chatID); err != nil {
		return errDBProblem
	}

	msg := tgbotapi.NewMessage(chatID, b.messages.Start)
	b.sendMessage(msg)
	return nil
}

func (b *Bot) handleHelpCommand(chatID int64) error {
	msg := tgbotapi.NewMessage(chatID, b.messages.Help)
	b.sendMessage(msg)
	return nil
}

func (b *Bot) handleSettingsCommand(chatID int64) error {

	err, settings := b.settingsRepository.Get(chatID)
	if err != nil {
		return errDBProblem
	}

	text := fmt.Sprintf(b.messages.Settings, settings.SecondaryLang, settings.Translations, settings.Examples, settings.Synonyms, settings.Sentences)
	text = strings.ReplaceAll(text, "true", "on")
	text = strings.ReplaceAll(text, "false", "off")

	msg := tgbotapi.NewMessage(chatID, text)
	b.sendMessage(msg)
	return nil
}

func (b *Bot) handleSetSecondaryLangCommand(chatID int64) error {
	if err := b.statesRepository.Save(chatID, chooseLanguageState); err != nil {
		return errDBProblem
	}

	msg := tgbotapi.NewMessage(chatID, b.messages.ChooseLang)
	msg.ReplyMarkup = langKeyboard()
	b.sendMessage(msg)
	return nil
}

func (b *Bot) handleToggleSomethingCommand(chatID int64, toggle func(s *service.Settings) bool) error {
	err, settings := b.settingsRepository.Get(chatID)
	if err != nil {
		return errDBProblem
	}
	status := toggle(settings)
	if err := b.settingsRepository.Save(chatID, settings); err != nil {
		return errDBProblem
	}
	text := b.messages.SettingOn
	if !status {
		text = b.messages.SettingOff
	}
	msg := tgbotapi.NewMessage(chatID, text)
	b.sendMessage(msg)
	return nil
}
