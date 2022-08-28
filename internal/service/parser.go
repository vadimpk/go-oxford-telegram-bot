package service

import (
	"github.com/vadimpk/go-oxford-dictionary-sdk"
	"strconv"
)

type OxfordParser struct {
	oxfordClient *oxford.Client
}

func NewOxfordParser(oxfordClient *oxford.Client) *OxfordParser {
	return &OxfordParser{oxfordClient: oxfordClient}
}

func (p *OxfordParser) Parse(word string, settings *Settings) (error, string) {
	entry, err := p.oxfordClient.Entry(word)
	if err != nil {
		return err, ""
	}

	var translations []string
	var sentences []string

	if settings.Translations {
		resp, err := p.oxfordClient.Translation(word, "en", settings.SecondaryLang)
		if err != nil {
			return err, ""
		}
		translations = parseTranslations(resp)
	}

	if settings.Sentences {
		resp, err := p.oxfordClient.Sentences(word)
		if err != nil {
			return err, ""
		}
		sentences = parseSentences(resp)
	}

	return p.print(entry, settings, translations, sentences)

}

func (p *OxfordParser) print(d oxford.OxfordResponse, settings *Settings, translations []string, sentences []string) (error, string) {

	var result string

	for _, le := range d.Results[0].LexicalEntries {
		result += le.Text + "\t" + le.Entries[0].Pronunciations[0].PhoneticSpelling
		result += le.LexicalCategory.Text + "\n"

		if len(translations) > 0 {
			for _, t := range translations {
				result += t + " "
			}
			result += "\n"
		}

		for sID, s := range le.Entries[0].Senses {
			result += "\n" + strconv.Itoa(sID) + ". " + s.Definitions[0] + "\n"

			if settings.Examples {
				for id, ex := range s.Examples {
					if id < 2 {
						result += ex.Text + "\n"
					}
				}
			}

			if settings.Synonyms {
				for id, s := range s.Synonyms {
					if id < 6 {
						result += s.Text + " "
					}
				}
				result += "\n"
			}
		}
	}

	if len(sentences) > 0 {
		result += "\n"
		for _, s := range sentences {
			result += s + "\n"
		}
	}

	return nil, result
}

func parseTranslations(d oxford.OxfordResponse) []string {
	var translations []string

	for _, le := range d.Results[0].LexicalEntries {
		for _, s := range le.Entries[0].Senses {
			for _, t := range s.Translations {
				translations = append(translations, t.Text)
			}
		}
	}
	return translations
}

func parseSentences(d oxford.OxfordResponse) []string {
	var sentences []string

	for _, le := range d.Results[0].LexicalEntries {
		for _, s := range le.Sentences {
			sentences = append(sentences, s.Text)
		}
	}
	return sentences
}
