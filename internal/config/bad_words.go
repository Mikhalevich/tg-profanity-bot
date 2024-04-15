package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed bad_words.json
var badWordsData []byte

func BadWords() ([]string, error) {
	var words []string
	if err := json.Unmarshal(badWordsData, &words); err != nil {
		return nil, fmt.Errorf("json unmarshal words: %w", err)
	}

	return words, nil
}
