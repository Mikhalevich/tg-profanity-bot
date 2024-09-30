package port

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type FormatType string

const (
	FormatTypeBold    FormatType = "bold"
	FormatTypeMention FormatType = "text_mention"
)

type Options struct {
	Buttons []*Button
	Format  []Format
}

type Option func(*Options)

type Button struct {
	Caption string
	Data    string
}

type Format struct {
	Type   FormatType
	Offset int
	Length int
	// used only for mention type
	User *tgbotapi.User
}

func WithButton(button *Button) Option {
	return func(opts *Options) {
		opts.Buttons = append(opts.Buttons, button)
	}
}

func WithFormat(format Format) Option {
	return func(opts *Options) {
		opts.Format = append(opts.Format, format)
	}
}

type MsgSender interface {
	Reply(ctx context.Context, originMsgInfo MessageInfo, msg string, options ...Option) error
	Edit(ctx context.Context, originMsgInfo MessageInfo, msg string, options ...Option) error
}
