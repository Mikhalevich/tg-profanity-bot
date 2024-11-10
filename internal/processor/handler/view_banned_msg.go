package handler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (h *handler) ViewBannedMsgCallbackQuery(
	ctx context.Context,
	info port.MessageInfo,
	msgText string,
) error {
	mangledMsg, err := h.mangler.Mangle(ctx, info.ChatID.String(), msgText)
	if err != nil {
		return fmt.Errorf("replace msg: %w", err)
	}

	if mangledMsg == msgText {
		if err := h.msgSender.Reply(ctx, info, msgText); err != nil {
			return fmt.Errorf("origin reply: %w", err)
		}

		return nil
	}

	if err := h.msgSender.Reply(
		ctx,
		info,
		mangledMsg,
		port.WithButton(h.viewOriginMsgButton(ctx, msgText)),
	); err != nil {
		return fmt.Errorf("mangled reply: %w", err)
	}

	return nil
}
