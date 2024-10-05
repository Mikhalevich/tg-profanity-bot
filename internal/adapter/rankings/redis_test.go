package rankings

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

type RedisSuit struct {
	*suite.Suite
	client   redis.UniversalClient
	cleanup  func() error
	rankings *redisRankings
}

func TestRedisSuit(t *testing.T) {
	t.Parallel()
	suite.Run(t, &RedisSuit{
		Suite: new(suite.Suite),
	})
}

func (s *RedisSuit) SetupSuite() {
	client, cleanup, err := redisConnection()
	if err != nil {
		s.FailNow("connect to redis", err)
	}

	s.client = client
	s.cleanup = cleanup
	s.rankings = NewRedisRankings(client, time.Hour)
}

func (s *RedisSuit) TearDownSuite() {
	if err := s.cleanup(); err != nil {
		s.FailNow("cleanup", err)
	}
}

func (s *RedisSuit) TearDownTest() {
	if err := s.client.FlushDB(context.Background()).Err(); err != nil {
		s.FailNow("clear all keys", err)
	}
}

func redisConnection() (*redis.Client, func() error, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("construct pool: %w", err)
	}

	if err := pool.Client.Ping(); err != nil {
		return nil, nil, fmt.Errorf("connect to docker: %w", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "redis",
		Tag:        "7.0.15",
		Env: []string{
			"REDIS_PASSWORD=redis123",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})

	if err != nil {
		return nil, nil, fmt.Errorf("run docker: %w", err)
	}

	var client *redis.Client

	if err := pool.Retry(func() error {
		client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("localhost:%s", resource.GetPort("6379/tcp")),
			Password: "redis123",
			DB:       1,
		})

		if err := client.Ping(context.Background()).Err(); err != nil {
			return fmt.Errorf("redis ping: %w", err)
		}

		return nil
	}); err != nil {
		return nil, nil, fmt.Errorf("connect to database: %w", err)
	}

	return client, func() error {
		if err := pool.Purge(resource); err != nil {
			return fmt.Errorf("purge resource: %w", err)
		}

		return nil
	}, nil
}

func (s *RedisSuit) TestAddAndGetTopScores() {
	var (
		ctx = context.Background()
		key = "tankings_test_key"

		userID1 = "user_id_1"
		userID2 = "user_id_2"
		userID3 = "user_id_3"
	)

	s.Require().NoError(s.rankings.AddScore(ctx, key, userID1))
	s.Require().NoError(s.rankings.AddScore(ctx, key, userID1))
	s.Require().NoError(s.rankings.AddScore(ctx, key, userID1))

	s.Require().NoError(s.rankings.AddScore(ctx, key, userID2))

	s.Require().NoError(s.rankings.AddScore(ctx, key, userID3))
	s.Require().NoError(s.rankings.AddScore(ctx, key, userID3))

	top, err := s.rankings.Top(ctx, key)
	s.Require().NoError(err)

	s.Require().ElementsMatch(top, []port.RankingUserScore{
		{
			UserID: userID1,
			Score:  3,
		},
		{
			UserID: userID3,
			Score:  2,
		},
		{
			UserID: userID2,
			Score:  1,
		},
	})
}

func (s *RedisSuit) TestTopForEmptyUsers() {
	var (
		ctx = context.Background()
		key = "tankings_test_key"
	)

	top, err := s.rankings.Top(ctx, key)
	s.Require().NoError(err)

	s.Require().Empty(top)
}
