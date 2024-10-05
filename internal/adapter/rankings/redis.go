package rankings

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

var _ port.Rankings = (*redisRankings)(nil)

type redisRankings struct {
	client redis.UniversalClient
	ttl    time.Duration
}

func NewRedisRankings(client redis.UniversalClient, ttl time.Duration) *redisRankings {
	return &redisRankings{
		client: client,
		ttl:    ttl,
	}
}

func (r *redisRankings) AddScore(ctx context.Context, key string, userID string) error {
	if err := r.incrByAlreadyExistingKey(ctx, key, userID); err != nil {
		if !errors.Is(err, redis.Nil) {
			return fmt.Errorf("incr by already existing key: %w", err)
		}

		if err := r.incrByWithExpiration(ctx, key, userID); err != nil {
			return fmt.Errorf("incr by with expiration: %w", err)
		}

		return nil
	}

	return nil
}

func (r *redisRankings) incrByAlreadyExistingKey(ctx context.Context, key string, userID string) error {
	if err := r.client.ZAddArgsIncr(ctx, key, redis.ZAddArgs{
		XX: true,
		Members: []redis.Z{
			{
				Score:  1.0,
				Member: userID,
			},
		},
	}).Err(); err != nil {
		return fmt.Errorf("zaddargsincr: %w", err)
	}

	return nil
}

func (r *redisRankings) incrByWithExpiration(ctx context.Context, key string, userID string) error {
	if _, err := r.client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		if err := pipe.ZIncrBy(ctx, key, 1.0, userID).Err(); err != nil {
			return fmt.Errorf("zincrby: %w", err)
		}

		if err := pipe.Expire(ctx, key, r.ttl).Err(); err != nil {
			return fmt.Errorf("expire: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("pipelined: %w", err)
	}

	return nil
}

func (r *redisRankings) Top(ctx context.Context, key string) ([]port.RankingUserScore, error) {
	rawScores, err := r.client.ZRevRangeWithScores(ctx, key, 0, -1).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// key does not exists
			return nil, nil
		}

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
