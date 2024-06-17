package processor

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Mikhalevich/tg-profanity-bot/internal/app/tracing"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *processor) ProcessMessage(ctx context.Context, msg *tgbotapi.Message) error {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	mangledMsg, err := p.replacer.Replace(ctx, strconv.FormatInt(msg.Chat.ID, 10), msg.Text)
	if err != nil {
		return fmt.Errorf("replace msg: %w", err)
	}

	if mangledMsg == msg.Text {
		return nil
	}

	if err := p.msgSender.Edit(ctx, msg, mangledMsg); err != nil {
		return fmt.Errorf("msg edit: %w", err)
	}

	return nil
}
