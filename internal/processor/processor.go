package processor

import (
	"context"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/cmd"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TextReplacer interface {
	Replace(ctx context.Context, chatID string, msg string) (string, error)
}

type MsgSender interface {
	Reply(ctx context.Context, originMsg *tgbotapi.Message, msg string) error
	Edit(ctx context.Context, originMsg *tgbotapi.Message, msg string) error
}

type WordsProvider interface {
	ChatWords(ctx context.Context, chatID string) ([]string, error)
}

type processor struct {
	replacer      TextReplacer
	msgSender     MsgSender
	wordsProvider WordsProvider
	router        cmd.Router
}

func New(replacer TextReplacer, msgSender MsgSender, wordsProvider WordsProvider) *processor {
	p := &processor{
		replacer:      replacer,
		msgSender:     msgSender,
		wordsProvider: wordsProvider,
	}

	p.initRoutes()

	return p
}

func (p *processor) initRoutes() {
	p.router = cmd.Router{
		"getall": {
			Handler: p.GetAllWords,
			Perm:    cmd.Admin,
		},
	}
}
