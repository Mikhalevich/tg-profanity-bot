package banprocessor

import (
	"context"
)

type nope struct {
}

func NewNope() *nope {
	return &nope{}
}

func (n *nope) IsBanned(ctx context.Context, id string) bool {
	return false
}

func (n *nope) Unban(ctx context.Context, id string) error {
	return nil
}

func (n *nope) AddViolation(ctx context.Context, id string) (bool, error) {
	return false, nil
}
