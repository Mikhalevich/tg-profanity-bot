package matcher

import (
	"fmt"
)

type ChatWordsProvider interface {
	ChatWords(chatID string) ([]string, error)
	CreateChatWords(chatID string, words []string) error
	IsChatNotExistsError(err error) bool
}

type ahocorasickDynamicMatcher struct {
	cwp ChatWordsProvider
}

func NewNewAhocorasickDynamic(provider ChatWordsProvider) *ahocorasickDynamicMatcher {
	return &ahocorasickDynamicMatcher{
		cwp: provider,
	}
}

func (m *ahocorasickDynamicMatcher) Match(chatID string, in []byte) ([]string, error) {
	words, err := m.cwp.ChatWords(chatID)
	if err != nil {
		return nil, fmt.Errorf("chat words: %w", err)
	}

	matchedWords, err := NewAhocorasick(words).Match(chatID, in)
	if err != nil {
		return nil, fmt.Errorf("arhocarasick match: %w", err)
	}

	return matchedWords, nil
}
