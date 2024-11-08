package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Mikhalevich/tg-profanity-bot/internal/app/messagequeue/rabbit/consumer"
	"github.com/Mikhalevich/tg-profanity-bot/internal/config"
	"github.com/Mikhalevich/tg-profanity-bot/internal/infra"
	"github.com/Mikhalevich/tg-profanity-bot/internal/infra/logger"
	"github.com/Mikhalevich/tg-profanity-bot/internal/infra/tracing"
)

func main() {
	var cfg config.Consumer
	if err := infra.LoadConfig(&cfg); err != nil {
		logger.StdLogger().WithError(err).Error("failed to load config")
		os.Exit(1)
	}

	l, err := infra.SetupLogger(cfg.LogLevel)
	if err != nil {
		logger.StdLogger().WithError(err).Error("failed to setup logger")
		os.Exit(1)
	}

	if err := runService(cfg, l); err != nil {
		l.WithError(err).Error("failed run service")
		os.Exit(1)
	}
}

func runService(cfg config.Consumer, l logger.Logger) error {
	if err := tracing.SetupTracer(cfg.Tracing.Endpoint, cfg.Tracing.ServiceName, ""); err != nil {
		return fmt.Errorf("setup tracer: %w", err)
	}

	msgProcessor, cleanup, err := infra.MakeMsgProcessor(
		cfg.BotToken,
		cfg.Postgres,
		cfg.Profanity,
		cfg.CommandRedis,
		cfg.BanRedis,
		cfg.RankingsRedis,
	)
	if err != nil {
		return fmt.Errorf("init msg processor: %w", err)
	}

	defer cleanup()

	ch, channelCleanup, err := infra.MakeRabbitAMQPChannel(cfg.Rabbit.URL)
	if err != nil {
		return fmt.Errorf("init rabbitmq channel: %w", err)
	}

	defer channelCleanup()

	c, err := consumer.New(ch, cfg.Rabbit.MsgQueue, l.WithField("bot_name", "bot_msg_worker"))
	if err != nil {
		return fmt.Errorf("init message consumer: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	l.Info("consumer running...")

	if err := c.Consume(ctx, cfg.Rabbit.WorkersCount, msgProcessor); err != nil {
		return fmt.Errorf("consume messages: %w", err)
	}

	l.Info("consumer stopped...")

	return nil
}
