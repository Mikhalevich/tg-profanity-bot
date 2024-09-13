package processor

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (p *processor) GetAllWords(ctx context.Context, info port.MessageInfo, cmdArgs string) error {
	words, err := p.wordsProvider.ChatWords(ctx, info.ChatID.String())
	if err != nil {
		return fmt.Errorf("get chat words: %w", err)
	}

	if err := p.msgSender.Reply(ctx, info, msgFromWords(words)); err != nil {
		return fmt.Errorf("msg reply: %w", err)
	}

	return nil
}

func msgFromWords(words []string) string {
	if len(words) == 0 {
		return "words are empty"
	}

	return strings.Join(words, "\n")
}
