package text_info

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Definition []struct {
	Word      string `json:"word"`
	Phonetics []struct {
		Audio     string `json:"audio"`
		SourceURL string `json:"sourceUrl,omitempty"`
		License   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"license,omitempty"`
		Text string `json:"text,omitempty"`
	} `json:"phonetics"`
	Meanings []struct {
		PartOfSpeech string `json:"partOfSpeech"`
		Definitions  []struct {
			Definition string        `json:"definition"`
			Synonyms   []interface{} `json:"synonyms"`
			Antonyms   []interface{} `json:"antonyms"`
		} `json:"definitions"`
		Synonyms []string      `json:"synonyms"`
		Antonyms []interface{} `json:"antonyms"`
	} `json:"meanings"`
	License struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"license"`
	SourceUrls []string `json:"sourceUrls"`
}

func getDefinition(word string) (string, error) {
	resp, err := http.Get("https://api.dictionaryapi.dev/api/v2/entries/en/" + word)

	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var body Definition
	err = json.Unmarshal(b, &body)
	if err != nil {
		return "", err
	}

	var count int
	var response string
	for _, meaning := range body[0].Meanings {
		for _, def := range meaning.Definitions {
			if count >= 3 {
				break
			}
			count++
			strCount := strconv.Itoa(count)
			response += strCount + ": " + def.Definition + "" + "\n\n"
		}
	}

	return response, nil
}
