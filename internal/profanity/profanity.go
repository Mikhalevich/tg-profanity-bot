package profanity

import (
	"context"
)

type profanity struct {
}

func New() *profanity {
	return &profanity{}
}

func (p *profanity) ReplaceMessage(ctx context.Context, msg string) (string, error) {
	return "", nil
}
