package matcher

import (
	"context"
	"fmt"
)

type ChatWordsProvider interface {
	ChatWords(ctx context.Context, chatID string) ([]string, error)
	CreateChatWords(ctx context.Context, chatID string, words []string) error
	IsChatNotExistsError(err error) bool
}

type ahocorasickDynamicMatcher struct {
	cwp              ChatWordsProvider
	initialChatWords []string
}

func NewNewAhocorasickDynamic(provider ChatWordsProvider, initialChatWords []string) *ahocorasickDynamicMatcher {
	return &ahocorasickDynamicMatcher{
		cwp:              provider,
		initialChatWords: initialChatWords,
	}
}

func (m *ahocorasickDynamicMatcher) Match(chatID string, in []byte) ([]string, error) {
	words, err := m.cwp.ChatWords(context.Background(), chatID)
	if err != nil {
		if !m.cwp.IsChatNotExistsError(err) {
			return nil, fmt.Errorf("chat words: %w", err)
		}

		if err := m.cwp.CreateChatWords(context.Background(), chatID, m.initialChatWords); err != nil {
			return nil, fmt.Errorf("create chat words: %w", err)
		}

		return m.initialChatWords, nil
	}

	matchedWords, err := NewAhocorasick(words).Match(chatID, in)
	if err != nil {
		return nil, fmt.Errorf("arhocarasick match: %w", err)
	}

	return matchedWords, nil
}
