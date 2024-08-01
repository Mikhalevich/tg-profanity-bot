package processor

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Mikhalevich/tg-profanity-bot/internal/app/tracing"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/cmd"
)

func (p *processor) ProcessCallbackQuery(ctx context.Context, query *tgbotapi.CallbackQuery) error {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	command, err := p.commandStorage.Get(ctx, query.Data)
	if err != nil {
		if !p.commandStorage.IsNotFoundError(err) {
			return fmt.Errorf("get command from store: %w", err)
		}

		if err := p.msgSender.Reply(ctx, query.Message, "command expired"); err != nil {
			return fmt.Errorf("send command expired: %w", err)
		}
	}

	r, ok := p.buttonsRouter[cmd.CMD(command.CMD)]
	if !ok {
		return fmt.Errorf("unsupported command %s", command.CMD)
	}

	if r.IsAdmin() && !p.permissionChecker.IsAdmin(ctx, query.Message.Chat.ID, query.From.ID) {
		if err := p.msgSender.Reply(ctx, query.Message, "need admin permission"); err != nil {
			return fmt.Errorf("send perm message: %w", err)
		}

		return nil
	}

	var (
		chatID  = strconv.FormatInt(query.Message.Chat.ID, 10)
		payload = string(command.Payload)
	)

	if err := r.Handler(ctx, chatID, payload, query.Message); err != nil {
		return fmt.Errorf("handle query %s: %w", command.CMD, err)
	}

	return nil
}
