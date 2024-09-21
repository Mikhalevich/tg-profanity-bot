package rankings

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

var _ port.Rankings = (*redisRankings)(nil)

type redisRankings struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisRankings(client *redis.Client, ttl time.Duration) *redisRankings {
	return &redisRankings{
		client: client,
		ttl:    ttl,
	}
}

func (r *redisRankings) AddScore(ctx context.Context, key string, userID string) error {
	if err := r.client.ZIncrBy(ctx, key, 1.0, userID).Err(); err != nil {
		return fmt.Errorf("zincrby: %w", err)
	}

	return nil
}

func (r *redisRankings) Top(ctx context.Context, key string) ([]port.RankingUserScore, error) {
	rawScores, err := r.client.ZRevRangeWithScores(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("zrevrangewithscores: %w", err)
	}

	if len(rawScores) == 0 {
		return nil, nil
	}

	userScores := make([]port.RankingUserScore, 0, len(rawScores))

	for _, rw := range rawScores {
		userID, ok := rw.Member.(string)
		if !ok {
			return nil, fmt.Errorf("invalid user_id type: %v", rw.Member)
		}

		userScores = append(userScores, port.RankingUserScore{
			UserID: userID,
			Score:  int(rw.Score),
		})
	}

	return userScores, nil
}
