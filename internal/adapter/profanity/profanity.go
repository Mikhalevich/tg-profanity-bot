package profanity

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/profanity/internal/position"
	"github.com/Mikhalevich/tg-profanity-bot/internal/app/tracing"
)

type Matcher interface {
	Match(ctx context.Context, chatID string, in []byte) ([]string, error)
}

type Replacer interface {
	Replace(text string) string
}

type profanity struct {
	matcher  Matcher
	replacer Replacer
}

func New(matcher Matcher, replacer Replacer) *profanity {
	return &profanity{
		matcher:  matcher,
		replacer: replacer,
	}
}

func (p *profanity) Replace(
	ctx context.Context,
	chatID string,
	msg string,
) (string, error) {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	msgLower := strings.ToLower(msg)

	foundedWords, err := p.matcher.Match(ctx, chatID, []byte(msgLower))
	if err != nil {
		return "", fmt.Errorf("match words: %w", err)
	}

	if len(foundedWords) == 0 {
		return msg, nil
	}

	var (
		wordsPositions   = p.wordsPositions(ctx, msgLower, foundedWords)
		reducedPositions = p.reduceInnerPositions(ctx, wordsPositions)
	)

	if len(reducedPositions) == 0 {
		return msg, nil
	}

	return p.mangle(ctx, msg, reducedPositions), nil
}

func (p *profanity) wordsPositions(ctx context.Context, msg string, foundedWords []string) *position.SortedPositions {
	_, span := tracing.StartSpan(ctx)
	defer span.End()

	var (
		msgLen    = len(msg)
		positions = position.NewSortedPositions()
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

func (p *profanity) reduceInnerPositions(ctx context.Context, wordsPositions *position.SortedPositions) []int {
	_, span := tracing.StartSpan(ctx)
	defer span.End()

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

func (p *profanity) mangle(ctx context.Context, msg string, positions []int) string {
	_, span := tracing.StartSpan(ctx)
	defer span.End()

	var (
		builder   = strings.Builder{}
		lastIndex = 0
	)

	builder.Grow(len(msg))

	for i := 0; i < len(positions); i += 2 {
		builder.WriteString(msg[lastIndex:positions[i]])

		censoredText := p.replacer.Replace(msg[positions[i]:positions[i+1]])
		builder.WriteString(censoredText)

		lastIndex = positions[i+1]
	}

	builder.WriteString(msg[lastIndex:])

	return builder.String()
}
