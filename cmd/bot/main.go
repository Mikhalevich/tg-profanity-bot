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
	"github.com/Mikhalevich/tg-profanity-bot/internal/bot"
	"github.com/Mikhalevich/tg-profanity-bot/internal/config"
	"github.com/Mikhalevich/tg-profanity-bot/internal/messagequeue/rabbit/publisher"
)

func main() {
	var cfg config.Bot
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

func runService(cfg config.Bot, logger *logrus.Logger) error {
	if err := tracing.SetupTracer(cfg.Tracing.Endpoint, cfg.Tracing.ServiceName, ""); err != nil {
		return fmt.Errorf("setup tracer: %w", err)
	}

	msgProcessor, cleanup, err := makeProcessor(cfg.Rabbit, cfg.Postgres, cfg.Profanity, cfg.Updates.Token)
	if err != nil {
		return fmt.Errorf("init processor: %w", err)
	}

	defer cleanup()

	tgBot, err := bot.New(cfg.Updates.Token, msgProcessor, logger.WithField("bot_name", "profanity_bot"))
	if err != nil {
		return fmt.Errorf("init bot: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()

		logger.Info("bot running...")
		tgBot.ProcessUpdates(cfg.Updates.UpdateTimeoutSeconds)
	}()

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)

loop:
	for {
		select {
		case <-terminate:
			signal.Stop(terminate)
			logger.Info("stopping bot...")
			tgBot.Stop()
		case <-ctx.Done():
			break loop
		}
	}

	logger.Info("bot stopped...")

	return nil
}

func makeProcessor(
	rabbitCfg config.RabbitMQProducer,
	postgresCfg config.Postgres,
	profanityCfg config.Profanity,
	botToken string,
) (bot.MessageProcessor, func(), error) {
	if rabbitCfg.URL != "" {
		msgPublisher, cleanup, err := makeRabbitPublisher(rabbitCfg)
		if err != nil {
			return nil, nil, fmt.Errorf("rabbit publisher: %w", err)
		}

		return msgPublisher, cleanup, nil
	}

	pg, pgCleanup, err := app.InitPostgres(postgresCfg)
	if err != nil {
		return nil, nil, fmt.Errorf("init postgres: %w", err)
	}

	msgProcessor, err := app.MakeMsgProcessor(profanityCfg, botToken, pg)
	if err != nil {
		return nil, nil, fmt.Errorf("msg processor: %w", err)
	}

	return msgProcessor, pgCleanup, nil
}

func makeRabbitPublisher(rabbitCfg config.RabbitMQProducer) (bot.MessageProcessor, func(), error) {
	ch, cleanup, err := app.MakeRabbitAMQPChannel(rabbitCfg.URL)
	if err != nil {
		return nil, nil, fmt.Errorf("rabbit channel: %w", err)
	}

	msgPublisher, err := publisher.New(tracing.WrapChannel(ch), rabbitCfg.MsgQueue)
	if err != nil {
		return nil, nil, fmt.Errorf("rabbit publisher: %w", err)
	}

	return msgPublisher, cleanup, nil
}
