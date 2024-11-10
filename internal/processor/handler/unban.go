package handler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (h *handler) UnbanCallbackQuery(
	ctx context.Context,
	info port.MessageInfo,
	userID string,
) error {
	if err := h.banProcessor.Unban(ctx, makeBanID(info.ChatID.String(), userID)); err != nil {
		return fmt.Errorf("unban: %w", err)
	}

	if err := h.msgSender.Reply(ctx, info, "user unbanned successfully"); err != nil {
		return fmt.Errorf("success unban reply: %w", err)
	}

	return nil
}
