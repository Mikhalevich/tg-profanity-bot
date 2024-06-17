package processor

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TextReplacer interface {
	Replace(ctx context.Context, chatID string, msg string) (string, error)
}

type MsgSender interface {
	Edit(ctx context.Context, originMsg *tgbotapi.Message, msg string) error
}

type processor struct {
	replacer  TextReplacer
	msgSender MsgSender
}

func New(replacer TextReplacer, msgSender MsgSender) *processor {
	return &processor{
		replacer:  replacer,
		msgSender: msgSender,
	}
}
