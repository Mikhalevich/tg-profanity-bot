package processor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (p *processor) UnbanCallbackQuery(
	ctx context.Context,
	info port.MessageInfo,
	userID string,
) error {
	if err := p.banProcessor.Unban(ctx, makeBanID(info.ChatID.String(), userID)); err != nil {
		return fmt.Errorf("unban: %w", err)
	}

	if err := p.msgSender.Reply(ctx, info, "user unbanned successfully"); err != nil {
		return fmt.Errorf("success unban reply: %w", err)
	}

	return nil
}
