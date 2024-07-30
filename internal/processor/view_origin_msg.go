package processor

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *processor) ViewOriginMsgCallbackQuery(
	ctx context.Context,
	chatID string,
	originMsgText string,
	msg *tgbotapi.Message,
) error {
	if err := p.msgSender.Reply(ctx, msg, originMsgText); err != nil {
		return fmt.Errorf("view reply: %w", err)
	}

	return nil
}
