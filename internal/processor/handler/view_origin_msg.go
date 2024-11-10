package handler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (h *handler) ViewOriginMsgCallbackQuery(
	ctx context.Context,
	info port.MessageInfo,
	originMsgText string,
) error {
	if err := h.msgSender.Reply(ctx, info, originMsgText); err != nil {
		return fmt.Errorf("view reply: %w", err)
	}

	return nil
}
