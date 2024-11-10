package handler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (h *handler) ProcessBannedMessage(ctx context.Context, info port.MessageInfo) (bool, error) {
	if !h.banProcessor.IsBanned(ctx, makeBanID(info.ChatID.String(), info.UserID.String())) {
		return false, nil
	}

	if err := h.editBanMessage(ctx, info.UserID.String(), info); err != nil {
		return false, fmt.Errorf("process ban: %w", err)
	}

	return true, nil
}

func makeBanID(chatID, userID string) string {
	return fmt.Sprintf("%s:%s", chatID, userID)
}

func (h *handler) editBanMessage(ctx context.Context, userID string, info port.MessageInfo) error {
	if err := h.msgSender.Edit(
		ctx,
		info,
		"user is banned",
		port.WithButton(h.viewBannedMsgButton(ctx, info.Text)),
		port.WithButton(h.unbanButton(ctx, userID)),
	); err != nil {
		return fmt.Errorf("msg edit: %w", err)
	}

	return nil
}
