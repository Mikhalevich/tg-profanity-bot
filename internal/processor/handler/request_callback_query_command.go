package handler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (h *handler) RequestCallbackQueryCommand(
	ctx context.Context,
	info port.MessageInfo,
	id string,
) (port.Command, bool, error) {
	command, err := h.commandStorage.Get(ctx, id)
	if err != nil {
		if !h.commandStorage.IsNotFoundError(err) {
			return port.Command{}, false, fmt.Errorf("get command from store: %w", err)
		}

		if err := h.msgSender.Reply(ctx, info, "command expired"); err != nil {
			return port.Command{}, false, fmt.Errorf("send command expired: %w", err)
		}
	}

	return command, false, nil
}
