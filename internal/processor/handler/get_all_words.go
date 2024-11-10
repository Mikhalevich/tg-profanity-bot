package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (h *handler) GetAllWords(ctx context.Context, info port.MessageInfo, cmdArgs string) error {
	words, err := h.wordsProvider.ChatWords(ctx, info.ChatID.String())
	if err != nil {
		return fmt.Errorf("get chat words: %w", err)
	}

	if err := h.msgSender.Reply(ctx, info, msgFromWords(words)); err != nil {
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
