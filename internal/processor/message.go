package processor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-profanity-bot/internal/app/tracing"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *processor) ProcessMessage(ctx context.Context, msg *tgbotapi.Message) error {
	ctx, span := tracing.Trace(ctx)
	defer span.End()

	mangledMsg := p.replacer.Replace(ctx, msg.Text)

	if mangledMsg == msg.Text {
		return nil
	}

	if err := p.msgSender.Edit(ctx, msg, mangledMsg); err != nil {
		return fmt.Errorf("msg edit: %w", err)
	}

	return nil
}
