package service

type Settings struct {
	SecondaryLang string
	Translations  bool
	Sentences     bool
	Synonyms      bool
	Examples      bool
}

func NewSettings() *Settings {
	return &Settings{SecondaryLang: "en", Translations: false, Sentences: false, Synonyms: false, Examples: false}
}

func (s *Settings) SetSecondaryLang(sl string) error {
	s.SecondaryLang = sl
	// TODO: check if exists
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
