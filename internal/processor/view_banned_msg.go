package processor

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *processor) ViewBannedMsgCallbackQuery(
	ctx context.Context,
	chatID string,
	msgText string,
	banMsg *tgbotapi.Message,
) error {
	mangledMsg, err := p.replacer.Replace(ctx, chatID, msgText)
	if err != nil {
		return fmt.Errorf("replace msg: %w", err)
	}

	if mangledMsg == msgText {
		if err := p.msgSender.Reply(ctx, banMsg, msgText); err != nil {
			return fmt.Errorf("origin reply: %w", err)
		}

		return nil
	}

	if err := p.msgSender.Reply(
		ctx,
		banMsg,
		mangledMsg,
		buttonRow(
			p.viewOriginMsgButton(ctx, msgText),
		)...,
	); err != nil {
		return fmt.Errorf("mangled reply: %w", err)
	}

	return nil
}
