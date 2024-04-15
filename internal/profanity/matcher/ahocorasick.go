package matcher

import (
	"github.com/cloudflare/ahocorasick"
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

func (am *ahocorasickMatcher) Match(in []byte) []string {
	wordsIdx := am.m.Match(in)
	if len(wordsIdx) == 0 {
		return nil
	}

	foundedWords := make([]string, 0, len(wordsIdx))
	for _, idx := range wordsIdx {
		foundedWords = append(foundedWords, am.words[idx])
	}

	return foundedWords
}
