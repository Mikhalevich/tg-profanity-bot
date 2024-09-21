package port

import (
	"context"
)

type RankingUser struct {
	ID          string
	DisplayName string
}

type RankingUserScore struct {
	User  RankingUser
	Score int
}

type Rankings interface {
	AddScore(ctx context.Context, key string, userInfo RankingUser) error
	Top(ctx context.Context, key string) ([]RankingUserScore, error)
}
