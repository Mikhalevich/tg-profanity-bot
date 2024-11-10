package handler

import (
	"context"

	"github.com/google/uuid"

	"github.com/Mikhalevich/tg-profanity-bot/internal/infra/logger"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/cbquery"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (h *handler) makeButton(
	ctx context.Context,
	caption string,
	command port.Command,
) *port.Button {
	id := uuid.NewString()

	if err := h.commandStorage.Set(ctx, id, command); err != nil {
		logger.FromContext(ctx).
			WithError(err).
			Warn("setting command storage error")

		return nil
	}

	return &port.Button{
		Caption: caption,
		Data:    id,
	}
}

func (h *handler) revertButton(ctx context.Context, button cbquery.CBQUERY, word string) *port.Button {
	return h.makeButton(ctx, "revert", port.Command{
		CMD:     button.String(),
		Payload: []byte(word),
	})
}

func (h *handler) viewOriginMsgButton(ctx context.Context, msg string) *port.Button {
	return h.makeButton(ctx, "view origin msg", port.Command{
		CMD:     cbquery.ViewOrginMsg.String(),
		Payload: []byte(msg),
	})
}

func (h *handler) viewBannedMsgButton(ctx context.Context, msg string) *port.Button {
	return h.makeButton(ctx, "view origin msg", port.Command{
		CMD:     cbquery.ViewBannedMsg.String(),
		Payload: []byte(msg),
	})
}

func (h *handler) unbanButton(ctx context.Context, userID string) *port.Button {
	return h.makeButton(ctx, "unban user", port.Command{
		CMD:     cbquery.Unban.String(),
		Payload: []byte(userID),
	})
}
