package processor

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TextReplacer interface {
	Replace(msg string) string
}

type MsgSender interface {
	Edit(originMsg *tgbotapi.Message, msg string) error
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
