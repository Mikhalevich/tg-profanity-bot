package matcher

import (
	"context"

	"github.com/cloudflare/ahocorasick"

	"github.com/Mikhalevich/tg-profanity-bot/internal/app/tracing"
)

type ahocorasickMatcher struct {
	m     ahocorasick.Matcher
	words []string
}

func NewAhocorasick(words []string) *ahocorasickMatcher {
	return &ahocorasickMatcher{
		m:     *ahocorasick.NewStringMatcher(words),
		words: words,
	}
}

func (am *ahocorasickMatcher) Match(ctx context.Context, chatID string, in []byte) ([]string, error) {
	_, span := tracing.StartSpan(ctx)
	defer span.End()

	wordsIdx := am.m.Match(in)
	if len(wordsIdx) == 0 {
		return nil, nil
	}

	foundedWords := make([]string, 0, len(wordsIdx))
	for _, idx := range wordsIdx {
		foundedWords = append(foundedWords, am.words[idx])
	}

	return foundedWords, nil
}
