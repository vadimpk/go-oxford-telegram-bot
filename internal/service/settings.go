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
	return &Settings{SecondaryLang: "es", Translations: false, Sentences: true, Synonyms: true, Examples: true}
}

func (s *Settings) SetSecondaryLang(sl string) error {
	if Languages[sl] == "" {
		return errors.New("invalid language")
	}
	s.SecondaryLang = Languages[sl]
	return nil
}

func ToggleTranslations(s *Settings) bool {
	s.Translations = !s.Translations
	return s.Translations
}

func ToggleSentences(s *Settings) bool {
	s.Sentences = !s.Sentences
	return s.Sentences
}

func ToggleSynonyms(s *Settings) bool {
	s.Synonyms = !s.Synonyms
	return s.Synonyms
}

func ToggleExamples(s *Settings) bool {
	s.Examples = !s.Examples
	return s.Examples
}
