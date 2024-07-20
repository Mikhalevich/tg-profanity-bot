package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/Mikhalevich/tg-profanity-bot/internal/app"
	"github.com/Mikhalevich/tg-profanity-bot/internal/app/tracing"
	"github.com/Mikhalevich/tg-profanity-bot/internal/config"
	"github.com/Mikhalevich/tg-profanity-bot/internal/messagequeue/rabbit/consumer"
)

func main() {
	var cfg config.Consumer
	if err := app.LoadConfig(&cfg); err != nil {
		logrus.WithError(err).Error("failed to load config")
		os.Exit(1)
	}

	logger, err := app.SetupLogger(cfg.LogLevel)
	if err != nil {
		logrus.WithError(err).Error("failed to setup logger")
		os.Exit(1)
	}

	if err := runService(cfg, logger); err != nil {
		logger.WithError(err).Error("failed run service")
		os.Exit(1)
	}
}

func runService(cfg config.Consumer, logger *logrus.Logger) error {
	if err := tracing.SetupTracer(cfg.Tracing.Endpoint, cfg.Tracing.ServiceName, ""); err != nil {
		return fmt.Errorf("setup tracer: %w", err)
	}

	msgProcessor, cleanup, err := app.MakeMsgProcessor(cfg.BotToken, cfg.Postgres, cfg.Profanity)
	if err != nil {
		return fmt.Errorf("init msg processor: %w", err)
	}

	defer cleanup()

	ch, channelCleanup, err := app.MakeRabbitAMQPChannel(cfg.Rabbit.URL)
	if err != nil {
		return fmt.Errorf("init rabbitmq channel: %w", err)
	}

	defer channelCleanup()

	c, err := consumer.New(ch, cfg.Rabbit.MsgQueue, logger.WithField("bot_name", "bot_msg_worker"))
	if err != nil {
		return fmt.Errorf("init message consumer: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger.Info("consumer running...")

	if err := c.Consume(ctx, cfg.Rabbit.WorkersCount, msgProcessor); err != nil {
		return fmt.Errorf("consume messages: %w", err)
	}

	logger.Info("consumer stopped...")

	return nil
}
