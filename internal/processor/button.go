package processor

import (
	"context"

	"github.com/google/uuid"

	"github.com/Mikhalevich/tg-profanity-bot/internal/app/logger"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/cbquery"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (p *processor) makeButton(
	ctx context.Context,
	caption string,
	command port.Command,
) *port.Button {
	id := uuid.NewString()

	if err := p.commandStorage.Set(ctx, id, command); err != nil {
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

func (p *processor) revertButton(ctx context.Context, button cbquery.CBQUERY, word string) *port.Button {
	return p.makeButton(ctx, "revert", port.Command{
		CMD:     button.String(),
		Payload: []byte(word),
	})
}

func (p *processor) viewOriginMsgButton(ctx context.Context, msg string) *port.Button {
	return p.makeButton(ctx, "view origin msg", port.Command{
		CMD:     cbquery.ViewOrginMsg.String(),
		Payload: []byte(msg),
	})
}

func (p *processor) viewBannedMsgButton(ctx context.Context, msg string) *port.Button {
	return p.makeButton(ctx, "view origin msg", port.Command{
		CMD:     cbquery.ViewBannedMsg.String(),
		Payload: []byte(msg),
	})
}

func (p *processor) unbanButton(ctx context.Context, userID string) *port.Button {
	return p.makeButton(ctx, "unban user", port.Command{
		CMD:     cbquery.Unban.String(),
		Payload: []byte(userID),
	})
}
