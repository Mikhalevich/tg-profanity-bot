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
	Reply(ctx context.Context, originMsg *tgbotapi.Message, msg string, buttons ...tgbotapi.InlineKeyboardButton) error
	Edit(ctx context.Context, originMsg *tgbotapi.Message, msg string, buttons ...tgbotapi.InlineKeyboardButton) error
}

type WordsProvider interface {
	ChatWords(ctx context.Context, chatID string) ([]string, error)
}

type WordsUpdater interface {
	AddWord(ctx context.Context, chatID string, word string) error
	RemoveWord(ctx context.Context, chatID string, word string) error
	IsNothingUpdatedError(err error) bool
}

type PermissionChecker interface {
	IsAdmin(ctx context.Context, chatID, userID int64) bool
}

type Command struct {
	CMD     string
	Payload []byte
}

type CommandStorage interface {
	Set(ctx context.Context, id string, command Command) error
	Get(ctx context.Context, id string) (Command, error)
	IsNotFoundError(err error) bool
}

type processor struct {
	replacer          TextReplacer
	msgSender         MsgSender
	wordsProvider     WordsProvider
	wordsUpdater      WordsUpdater
	permissionChecker PermissionChecker
	commandStorage    CommandStorage

	cmdRouter     cmd.Router
	buttonsRouter cmd.Router
}

func New(
	replacer TextReplacer,
	msgSender MsgSender,
	wordsProvider WordsProvider,
	wordsUpdater WordsUpdater,
	permissionChecker PermissionChecker,
	commandStorage CommandStorage,
) *processor {
	p := &processor{
		replacer:          replacer,
		msgSender:         msgSender,
		wordsProvider:     wordsProvider,
		wordsUpdater:      wordsUpdater,
		permissionChecker: permissionChecker,
		commandStorage:    commandStorage,
	}

	p.initCommandRoutes()
	p.initButtonsRoutes()

	return p
}

func (p *processor) initCommandRoutes() {
	p.cmdRouter = cmd.Router{
		cmd.GetAll: {
			Handler: p.GetAllWords,
			Perm:    cmd.Admin,
		},
	}

	if p.wordsUpdater != nil {
		p.cmdRouter[cmd.Add] = cmd.Route{
			Handler: p.AddWordCommand,
			Perm:    cmd.Admin,
		}

		p.cmdRouter[cmd.Remove] = cmd.Route{
			Handler: p.RemoveWordCommand,
			Perm:    cmd.Admin,
		}
	}
}

func (p *processor) initButtonsRoutes() {
	if p.wordsUpdater == nil {
		return
	}

	p.buttonsRouter = cmd.Router{
		cmd.Add: {
			Handler: p.AddWordCallbackQuery,
			Perm:    cmd.Admin,
		},
		cmd.Remove: {
			Handler: p.RemoveWordCallbackQuery,
			Perm:    cmd.Admin,
		},
		cmd.ViewOrginMsg: {
			Handler: p.ViewOriginMsgCallbackQuery,
			Perm:    cmd.Admin,
		},
	}
}
