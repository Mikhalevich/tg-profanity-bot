package commandstorage

import (
	"context"
	"errors"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor"
)

var (
	errNotFound = errors.New("not found")
)

type nope struct {
}

func NewNope() *nope {
	return &nope{}
}

func (n *nope) Set(ctx context.Context, id string, command processor.Command) error {
	return errors.New("not implemented")
}

func (n *nope) Get(ctx context.Context, id string) (processor.Command, error) {
	return processor.Command{}, errNotFound
}

func (n *nope) IsNotFoundError(err error) bool {
	return errors.Is(err, errNotFound)
}
