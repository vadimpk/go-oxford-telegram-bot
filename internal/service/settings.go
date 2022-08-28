package service

import (
	"errors"
)

type Settings struct {
	SecondaryLang string
	Translations  bool
	Sentences     bool
	Synonyms      bool
	Examples      bool
}

var Languages = map[string]string{
	"Arabic":     "ar",
	"German":     "de",
	"Russian":    "ru",
	"Spanish":    "es",
	"Chinese":    "zh",
	"Greek":      "el",
	"Hindi":      "hi",
	"Italian":    "it",
	"Portuguese": "pt",
	"Turkmen":    "tk",
}

func NewSettings() *Settings {
	return &Settings{SecondaryLang: "en", Translations: false, Sentences: false, Synonyms: false, Examples: false}
}

func (s *Settings) SetSecondaryLang(sl string) error {
	if Languages[sl] == "" {
		return errors.New("invalid language")
	}
	s.SecondaryLang = Languages[sl]
	return nil
}

func ToggleTranslations(s *Settings) {
	s.Translations = !s.Translations
}

func ToggleSentences(s *Settings) {
	s.Sentences = !s.Sentences
}

func ToggleSynonyms(s *Settings) {
	s.Synonyms = !s.Synonyms
}

func ToggleExamples(s *Settings) {
	s.Examples = !s.Examples
}
