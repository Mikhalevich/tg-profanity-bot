package staticwords

import (
	"context"
	"fmt"
)

type ChatWordsProvider interface {
	ChatWords(ctx context.Context, chatID string) ([]string, error)
}

type staticWords struct {
	words    []string
	provider ChatWordsProvider
}

func New(words []string, provider ChatWordsProvider) *staticWords {
	return &staticWords{
		words:    words,
		provider: provider,
	}
}

func (s *staticWords) ChatWords(ctx context.Context, chatID string) ([]string, error) {
	if s.provider == nil {
		return s.words, nil
	}

	words, err := s.provider.ChatWords(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("provider get chat words: %w", err)
	}

	return words, nil
}

func (s *staticWords) InitialWords() []string {
	return s.words
}
