package msgprocessor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-profanity-bot/internal/infra/tracing"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/cbquery"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (m *msgprocessor) ProcessCallbackQuery(ctx context.Context, query port.CallbackQuery) error {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	command, alreadyProcessed, err := m.handler.RequestCallbackQueryCommand(ctx, query.MessageInfo, query.Data)
	if err != nil {
		return fmt.Errorf("request cbq command: %w", err)
	}

	if alreadyProcessed {
		return nil
	}

	r, ok := m.buttonsRouter.Route(cbquery.CBQUERY(command.CMD))
	if !ok {
		return fmt.Errorf("unsupported command %s", command.CMD)
	}

	if r.IsAdmin() && !m.permissionChecker.IsAdmin(ctx, query.ChatID.Int64(), query.UserID.Int64()) {
		if err := m.handler.ReplyTextMessage(ctx, query.MessageInfo, "need admin permission"); err != nil {
			return fmt.Errorf("send perm message: %w", err)
		}

		return nil
	}

	if err := r.Handler(ctx, query.MessageInfo, string(command.Payload)); err != nil {
		return fmt.Errorf("handle query %s: %w", command.CMD, err)
	}

	return nil
}
