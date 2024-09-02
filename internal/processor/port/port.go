package port

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Mangler interface {
	Mangle(ctx context.Context, chatID string, msg string) (string, error)
}

type MsgSender interface {
	Reply(ctx context.Context, originMsgInfo MessageInfo, msg string, buttons ...tgbotapi.InlineKeyboardButton) error
	Edit(ctx context.Context, originMsgInfo MessageInfo, msg string, buttons ...tgbotapi.InlineKeyboardButton) error
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

type BanProcessor interface {
	IsBanned(ctx context.Context, id string) bool
	Unban(ctx context.Context, id string) error
	AddViolation(ctx context.Context, id string) (bool, error)
}
