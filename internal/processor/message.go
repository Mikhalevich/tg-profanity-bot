package processor

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *processor) ProcessMessage(ctx context.Context, msg *tgbotapi.Message) error {
	mangledMsg := p.replacer.Replace(msg.Text)

	if mangledMsg == msg.Text {
		return nil
	}

	if err := p.msgSender.Edit(msg, mangledMsg); err != nil {
		return fmt.Errorf("msg edit: %w", err)
	}

	return nil
}
