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
	Reply(ctx context.Context, originMsg *tgbotapi.Message, msg string, buttons ...[]tgbotapi.InlineKeyboardButton) error
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

type ChatMemberChecker interface {
	IsAdmin(chatID, userID int64) bool
}

type processor struct {
	replacer      TextReplacer
	msgSender     MsgSender
	wordsProvider WordsProvider
	wordsUpdater  WordsUpdater
	memberChecker ChatMemberChecker
	cmdRouter     cmd.Router
	buttonsRouter cmd.Router
}

func New(
	replacer TextReplacer,
	msgSender MsgSender,
	wordsProvider WordsProvider,
	wordsUpdater WordsUpdater,
	memberChecker ChatMemberChecker,
) *processor {
	p := &processor{
		replacer:      replacer,
		msgSender:     msgSender,
		wordsProvider: wordsProvider,
		wordsUpdater:  wordsUpdater,
		memberChecker: memberChecker,
	}

	p.initCommandRoutes()
	p.initButtonsRoutes()

	return p
}

func (p *processor) initCommandRoutes() {
	p.cmdRouter = cmd.Router{
		"getall": {
			Handler: p.GetAllWords,
			Perm:    cmd.Admin,
		},
	}

	if p.wordsUpdater != nil {
		p.cmdRouter["add"] = cmd.Route{
			Handler: p.AddWord,
			Perm:    cmd.Admin,
		}

		p.cmdRouter["remove"] = cmd.Route{
			Handler: p.RemoveWord,
			Perm:    cmd.Admin,
		}
	}
}

func (p *processor) initButtonsRoutes() {
	if p.wordsUpdater == nil {
		return
	}

	p.buttonsRouter = cmd.Router{
		"add": {
			Handler: p.AddWord,
			Perm:    cmd.Admin,
		},
		"remove": {
			Handler: p.RemoveWord,
			Perm:    cmd.Admin,
		},
	}
}
