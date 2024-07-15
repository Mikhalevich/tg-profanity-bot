package processor

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *processor) tryProcessCommand(ctx context.Context, chatID string, msg *tgbotapi.Message) (bool, error) {
	cmd, _ := extractCommand(msg.Text)
	if cmd == "" {
		return false, nil
	}

	switch cmd {
	case "getall":
		if err := p.GetAllWords(ctx, chatID, msg); err != nil {
			return false, fmt.Errorf("view words: %w", err)
		}
	default:
		return false, nil
	}

	return true, nil
}

func extractCommand(msg string) (string, string) {
	if !strings.HasPrefix(msg, "/") {
		return "", ""
	}

	cmd, args, _ := strings.Cut(msg[1:], " ")

	return cmd, args
}

func (p *processor) GetAllWords(ctx context.Context, chatID string, msg *tgbotapi.Message) error {
	words, err := p.wordsProcessor.ChatWords(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get chat words: %w", err)
	}

	if err := p.msgSender.Reply(ctx, msg, strings.Join(words, "\n")); err != nil {
		return fmt.Errorf("msg reply: %w", err)
	}

	return nil
}