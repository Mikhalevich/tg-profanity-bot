package profanity

import (
	"context"
)

type profanity struct {
}

func New() *profanity {
	return &profanity{}
}

func (p *profanity) ProcessMessage(ctx context.Context, command string, arguments string) (string, error) {
	return "", nil
}
