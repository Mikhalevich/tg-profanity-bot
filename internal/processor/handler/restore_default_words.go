package handler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (h *handler) RestoreDefaultWordsCommand(ctx context.Context, info port.MessageInfo, cmdArgs string) error {
	if err := h.wordsUpdater.CreateChatWords(ctx, info.ChatID.String(), h.wordsProvider.InitialWords()); err != nil {
		return fmt.Errorf("create chat words: %w", err)
	}

	if err := h.msgSender.Reply(ctx, info, "words updated successfully"); err != nil {
		return fmt.Errorf("success reply: %w", err)
	}

	return nil
}
