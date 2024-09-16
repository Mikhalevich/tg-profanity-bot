package processor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-profanity-bot/internal/app/tracing"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/cbquery"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (p *processor) ProcessCallbackQuery(ctx context.Context, query port.CallbackQuery) error {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	command, err := p.commandStorage.Get(ctx, query.Data)
	if err != nil {
		if !p.commandStorage.IsNotFoundError(err) {
			return fmt.Errorf("get command from store: %w", err)
		}

		if err := p.msgSender.Reply(ctx, query.MessageInfo, "command expired"); err != nil {
			return fmt.Errorf("send command expired: %w", err)
		}
	}

	r, ok := p.buttonsRouter.Route(cbquery.CBQUERY(command.CMD))
	if !ok {
		return fmt.Errorf("unsupported command %s", command.CMD)
	}

	if r.IsAdmin() && !p.permissionChecker.IsAdmin(ctx, query.ChatID.Int64(), query.UserID.Int64()) {
		if err := p.msgSender.Reply(ctx, query.MessageInfo, "need admin permission"); err != nil {
			return fmt.Errorf("send perm message: %w", err)
		}

		return nil
	}

	if err := r.Handler(ctx, query.MessageInfo, string(command.Payload)); err != nil {
		return fmt.Errorf("handle query %s: %w", command.CMD, err)
	}

	return nil
}
