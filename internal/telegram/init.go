package telegram

import "github.com/vadimpk/go-oxford-telegram-bot/internal/service"

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
