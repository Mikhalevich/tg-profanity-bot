package profanity

import (
	"strings"
	"unicode/utf8"

	"github.com/Mikhalevich/tg-profanity-bot/internal/profanity/internal/position"
)

type Matcher interface {
	Match(in []byte) []string
}

type profanity struct {
	matcher Matcher
	symbol  byte
}

func New(matcher Matcher, symbol byte) *profanity {
	return &profanity{
		matcher: matcher,
		symbol:  symbol,
	}
}

func (p *profanity) ReplaceMessage(msg string) string {
	var (
		wordsPositions   = p.wordsPositions(strings.ToLower(msg))
		reducedPositions = p.reduceInnerPositions(wordsPositions)
	)

	return p.mangle(msg, reducedPositions)
}

func (p *profanity) wordsPositions(msg string) *position.SortedPositions {
	var (
		msgLen       = len(msg)
		foundedWords = p.matcher.Match([]byte(msg))
		positions    = position.NewSortedPositions()
	)

	for _, badWord := range foundedWords {
		offset := 0

		for {
			if offset >= msgLen {
				break
			}

			startIdx := strings.Index(msg[offset:], badWord)
			if startIdx < 0 {
				break
			}

			startIdx += offset
			endIdx := startIdx + len(badWord)
			offset = endIdx

			positions.Append(&position.Position{Pos: startIdx, IsEnd: false})
			positions.Append(&position.Position{Pos: endIdx, IsEnd: true})
		}
	}

	return positions
}

func (p *profanity) reduceInnerPositions(wordsPositions *position.SortedPositions) []int {
	var (
		reducedPositions []int
		opened           = 0
	)

	for _, pos := range wordsPositions.Positions() {
		if pos.IsEnd {
			opened--
			if opened == 0 {
				reducedPositions = append(reducedPositions, pos.Pos)
			}

			continue
		}

		if opened == 0 {
			reducedPositions = append(reducedPositions, pos.Pos)
		}

		opened++
	}

	return reducedPositions
}

func (p *profanity) mangle(msg string, positions []int) string {
	var (
		builder   = strings.Builder{}
		lastIndex = 0
	)

	for i := 0; i < len(positions); i += 2 {
		b := strings.Repeat(string(p.symbol), utf8.RuneCountInString(msg[positions[i]:positions[i+1]]))
		builder.WriteString(msg[lastIndex:positions[i]])
		builder.WriteString(b)

		lastIndex = positions[i+1]
	}

	builder.WriteString(msg[lastIndex:])

	return builder.String()
}
