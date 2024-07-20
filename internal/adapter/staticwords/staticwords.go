package staticwords

import (
	"context"
)

type staticWords struct {
	words []string
}

func New(words []string) *staticWords {
	return &staticWords{
		words: words,
	}
}

func (s *staticWords) ChatWords(ctx context.Context, chatID string) ([]string, error) {
	return s.words, nil
}
