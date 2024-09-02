package commandstorage

import (
	"context"
	"errors"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

var (
	errNotFound = errors.New("not found")
)

type nope struct {
}

func NewNope() *nope {
	return &nope{}
}

func (n *nope) Set(ctx context.Context, id string, command port.Command) error {
	return errors.New("not implemented")
}

func (n *nope) Get(ctx context.Context, id string) (port.Command, error) {
	return port.Command{}, errNotFound
}

func (n *nope) IsNotFoundError(err error) bool {
	return errors.Is(err, errNotFound)
}
