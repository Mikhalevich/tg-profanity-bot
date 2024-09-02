package commandstorage

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

type redisStorage struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedis(client *redis.Client, ttl time.Duration) *redisStorage {
	return &redisStorage{
		client: client,
		ttl:    ttl,
	}
}

func (r *redisStorage) Set(ctx context.Context, id string, command port.Command) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(command); err != nil {
		return fmt.Errorf("gob encode: %w", err)
	}

	if err := r.client.Set(ctx, id, buf.Bytes(), r.ttl).Err(); err != nil {
		return fmt.Errorf("redis set: %w", err)
	}

	return nil
}

func (r *redisStorage) Get(ctx context.Context, id string) (port.Command, error) {
	b, err := r.client.GetDel(ctx, id).Bytes()
	if err != nil {
		return port.Command{}, fmt.Errorf("redis getdel: %w", err)
	}

	var cmd port.Command
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&cmd); err != nil {
		return port.Command{}, fmt.Errorf("gob decode: %w", err)
	}

	return cmd, nil
}

func (r *redisStorage) IsNotFoundError(err error) bool {
	return errors.Is(err, redis.Nil)
}
