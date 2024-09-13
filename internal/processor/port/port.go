package port

import (
	"context"
)

type Mangler interface {
	Mangle(ctx context.Context, chatID string, msg string) (string, error)
}

type WordsProvider interface {
	ChatWords(ctx context.Context, chatID string) ([]string, error)
	InitialWords() []string
}

type WordsUpdater interface {
	AddWord(ctx context.Context, chatID string, word string) error
	RemoveWord(ctx context.Context, chatID string, word string) error
	ClearWords(ctx context.Context, chatID string) error
	CreateChatWords(ctx context.Context, chatID string, words []string) error
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
