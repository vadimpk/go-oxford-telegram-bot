package repository

import "github.com/vadimpk/go-oxford-telegram-bot/internal/service"

type SettingsRepository interface {
	Save(chatID int64, settings *service.Settings) error
	Get(chatID int64) (error, *service.Settings)
}

type StatesRepository interface {
	Save(chatID int64, state string) error
	Get(chatID int64) (error, string)
}
