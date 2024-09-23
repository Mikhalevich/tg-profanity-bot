package port

import (
	"context"
)

type RankingUserScore struct {
	UserID string
	Score  int
}

type Rankings interface {
	AddScore(ctx context.Context, key string, userID string) error
	Top(ctx context.Context, key string) ([]RankingUserScore, error)
}
