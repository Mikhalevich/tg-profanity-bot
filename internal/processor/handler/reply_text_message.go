package handler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (h *handler) ReplyTextMessage(ctx context.Context, originInfo port.MessageInfo, text string) error {
	if err := h.msgSender.Reply(ctx, originInfo, text); err != nil {
		return fmt.Errorf("send perm message: %w", err)
	}

	return nil
}
