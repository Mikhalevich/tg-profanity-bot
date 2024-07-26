package processor

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Mikhalevich/tg-profanity-bot/internal/app/tracing"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/button"
)

func (p *processor) ProcessCallbackQuery(ctx context.Context, query *tgbotapi.CallbackQuery) error {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	buttonInfo, err := button.FromBase64(query.Data)
	if err != nil {
		return fmt.Errorf("decode base64 data: %w", err)
	}

	r, ok := p.buttonsRouter[buttonInfo.CMD]
	if !ok {
		return fmt.Errorf("unsupported command %s", buttonInfo.CMD)
	}

	if r.IsAdmin() && !p.memberChecker.IsAdmin(query.Message.Chat.ID, query.From.ID) {
		if err := p.msgSender.Reply(ctx, query.Message, "need admin permission"); err != nil {
			return fmt.Errorf("send perm message: %w", err)
		}

		return nil
	}

	if err := r.Handler(ctx, strconv.FormatInt(query.Message.Chat.ID, 10), buttonInfo.Word, query.Message); err != nil {
		return fmt.Errorf("handle query %s: %w", buttonInfo.CMD, err)
	}

	return nil
}
