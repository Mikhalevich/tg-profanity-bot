package processor

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *processor) UnbanCallbackQuery(
	ctx context.Context,
	chatID string,
	userID string,
	msg *tgbotapi.Message,
) error {
	if err := p.banProcessor.Unban(ctx, makeBanID(chatID, userID)); err != nil {
		return fmt.Errorf("unban: %w", err)
	}

	if err := p.msgSender.Reply(ctx, msg, "user unbanned successfully"); err != nil {
		return fmt.Errorf("success unban reply: %w", err)
	}

	return nil
}
