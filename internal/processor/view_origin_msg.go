package processor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (p *processor) ViewOriginMsgCallbackQuery(
	ctx context.Context,
	info port.MessageInfo,
	originMsgText string,
) error {
	if err := p.msgSender.Reply(ctx, info, originMsgText); err != nil {
		return fmt.Errorf("view reply: %w", err)
	}

	return nil
}
