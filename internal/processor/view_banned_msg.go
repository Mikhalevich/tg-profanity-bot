package processor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (p *processor) ViewBannedMsgCallbackQuery(
	ctx context.Context,
	info port.MessageInfo,
	msgText string,
) error {
	mangledMsg, err := p.mangler.Mangle(ctx, info.ChatID.String(), msgText)
	if err != nil {
		return fmt.Errorf("replace msg: %w", err)
	}

	if mangledMsg == msgText {
		if err := p.msgSender.Reply(ctx, info, msgText); err != nil {
			return fmt.Errorf("origin reply: %w", err)
		}

		return nil
	}

	if err := p.msgSender.Reply(
		ctx,
		info,
		mangledMsg,
		p.viewOriginMsgButton(ctx, msgText),
	); err != nil {
		return fmt.Errorf("mangled reply: %w", err)
	}

	return nil
}
