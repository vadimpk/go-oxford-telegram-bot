package oxford

import (
	"fmt"
	"github.com/vadimpk/go-oxford-dictionary-sdk"
	"github.com/vadimpk/go-oxford-telegram-bot/internal/service"
	"strings"
)

type Parser struct {
	oxfordClient *oxford.Client
}

const (
	headerTemplate      = "*%s*   \\[%s]\n_%s_\n\n"
	translationTemplate = "_%s_, "
	definitionTemplate  = "%d. %s\n\n"
	exampleTemplate     = "%s_%s_\n"
	sentenceTemplate    = "• %s\n"
	tab                 = "      "
)

func NewOxfordParser(oxfordClient *oxford.Client) *Parser {
	return &Parser{oxfordClient: oxfordClient}
}

func (p *Parser) Parse(word string, settings *service.Settings) (error, string) {
	entry, err := p.oxfordClient.Entry(word)
	if err != nil {
		if err.Error() == "API Error: 404 Not Found" {
			return nil, ""
		}
		return err, ""
	}

	if len(entry.Results) == 0 {
		return nil, ""
	}

	if settings.Translations {
		resp, err := p.oxfordClient.Translation(word, "en", settings.SecondaryLang)
		if err != nil {
			return err, ""
		}
		addTranslations(&entry, resp)
	}

	if settings.Sentences {
		resp, err := p.oxfordClient.Sentences(word)
		if err != nil {
			return err, ""
		}
		addSentences(&entry, resp)
	}

	return p.print(entry, settings)

}

func (p *Parser) print(d oxford.OxfordResponse, settings *service.Settings) (error, string) {

	var result string

	for _, le := range d.Results[0].LexicalEntries {
		result += fmt.Sprintf(headerTemplate, le.Text, le.Entries[0].Pronunciations[0].PhoneticSpelling, strings.ToLower(le.LexicalCategory.Text))

		if len(le.Entries[0].Senses[0].Translations) > 0 {
			printTranslations(&result, le.Entries[0].Senses[0].Translations)
		}

		for sID, s := range le.Entries[0].Senses {
			if len(s.Definitions) == 0 {
				continue
			}
			result += fmt.Sprintf(definitionTemplate, sID+1, s.Definitions[0])

			if settings.Examples && len(s.Examples) > 0 {
				printExamples(&result, s.Examples)
			}

			if settings.Synonyms && len(s.Synonyms) > 0 {
				printSynonyms(&result, s.Synonyms)
			}
		}

		if settings.Sentences && len(le.Sentences) > 0 {
			printSentences(&result, le.Sentences)
		}
	}
	return nil, result
}

func printTranslations(result *string, tr []oxford.Translation) {
	for _, t := range tr {
		*result += fmt.Sprintf(translationTemplate, t.Text)
	}
	*result = strings.TrimRight(*result, ", ")
	*result += "\n\n"
}

func printExamples(result *string, e []oxford.Example) {
	for id, ex := range e {
		if id > 1 {
			break
		}
		*result += fmt.Sprintf(exampleTemplate, tab, ex.Text)
	}
	*result += "\n"
}

func printSynonyms(result *string, s []oxford.Synonym) {
	*result += tab + "_synonyms:_ "
	for id, s := range s {
		if id > 4 {
			break
		}
		*result += s.Text + ", "
	}
	*result = strings.TrimRight(*result, ", ")
	*result += "\n\n"
}

func printSentences(result *string, sn []oxford.Sentence) {
	*result += "_Sentences_:\n"
	for id, s := range sn {
		if id > 4 {
			break
		}
		*result += fmt.Sprintf(sentenceTemplate, s.Text)
	}
	*result += "\n\n"
}

func addTranslations(entry *oxford.OxfordResponse, tr oxford.OxfordResponse) {
	var translations []oxford.Translation
	for leID := range entry.Results[0].LexicalEntries {
		if len(tr.Results[0].LexicalEntries) > leID {
			for _, s := range tr.Results[0].LexicalEntries[leID].Entries[0].Senses {
				for _, t := range s.Translations {
					translations = append(translations, t)
				}
			}
			entry.Results[0].LexicalEntries[leID].Entries[0].Senses[0].Translations = translations
			translations = translations[:0]
		}
	}
}

func addSentences(entry *oxford.OxfordResponse, s oxford.OxfordResponse) {
	for leID := range entry.Results[0].LexicalEntries {
		if len(s.Results[0].LexicalEntries) > leID {
			entry.Results[0].LexicalEntries[leID].Sentences = s.Results[0].LexicalEntries[leID].Sentences
		}
	}
}
