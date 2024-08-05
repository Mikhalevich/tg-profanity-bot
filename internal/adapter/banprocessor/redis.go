package banprocessor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

type redisBanProcessor struct {
	client            *redis.Client
	banTTL            time.Duration
	violationsPerHour int
	rate              *redis_rate.Limiter
}

func NewRedisBanProcessor(client *redis.Client, banTTL time.Duration, violationsPerHour int) *redisBanProcessor {
	return &redisBanProcessor{
		client:            client,
		banTTL:            banTTL,
		violationsPerHour: violationsPerHour,
		rate:              redis_rate.NewLimiter(client),
	}
}

func (r *redisBanProcessor) IsBanned(ctx context.Context, id string) bool {
	if err := r.client.Exists(ctx, makeBanKey(id)).Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return false
		}

		// skip error

		return false
	}

	return true
}

func makeBanKey(id string) string {
	return fmt.Sprintf("ban:%s", id)
}

func makeRateLimitterKey(id string) string {
	return fmt.Sprintf("rate:%s", id)
}

func (r *redisBanProcessor) AddViolation(ctx context.Context, id string) error {
	key := makeRateLimitterKey(id)

	res, err := r.rate.Allow(ctx, key, redis_rate.PerHour(r.violationsPerHour))
	if err != nil {
		return fmt.Errorf("rate allow: %w", err)
	}

	if res.Allowed > 0 {
		return nil
	}

	if err := r.client.Set(ctx, makeBanKey(id), "banned", r.banTTL).Err(); err != nil {
		return fmt.Errorf("set ban key: %w", err)
	}

	if err := r.rate.Reset(ctx, key); err != nil {
		return fmt.Errorf("rate reset: %w", err)
	}

	return nil
}
