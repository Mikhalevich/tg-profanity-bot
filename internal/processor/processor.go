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

type WordsUpdater interface {
	AddWord(ctx context.Context, chatID string, word string) error
	RemoveWord(ctx context.Context, chatID string, word string) error
	IsNothingUpdatedError(err error) bool
}

type processor struct {
	replacer      TextReplacer
	msgSender     MsgSender
	wordsProvider WordsProvider
	wordsUpdater  WordsUpdater
	router        cmd.Router
}

func New(
	replacer TextReplacer,
	msgSender MsgSender,
	wordsProvider WordsProvider,
	wordsUpdater WordsUpdater,
) *processor {
	p := &processor{
		replacer:      replacer,
		msgSender:     msgSender,
		wordsProvider: wordsProvider,
		wordsUpdater:  wordsUpdater,
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

	if p.wordsUpdater != nil {
		p.router["add"] = cmd.Route{
			Handler: p.AddWord,
			Perm:    cmd.Admin,
		}

		p.router["remove"] = cmd.Route{
			Handler: p.RemoveWord,
			Perm:    cmd.Admin,
		}
	}
}
