package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Mikhalevich/tg-profanity-bot/internal/app/bot"
	"github.com/Mikhalevich/tg-profanity-bot/internal/app/messagequeue/rabbit/publisher"
	"github.com/Mikhalevich/tg-profanity-bot/internal/config"
	"github.com/Mikhalevich/tg-profanity-bot/internal/infra"
	"github.com/Mikhalevich/tg-profanity-bot/internal/infra/logger"
	"github.com/Mikhalevich/tg-profanity-bot/internal/infra/tracing"
)

func main() {
	var cfg config.Bot
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

func runService(cfg config.Bot, l logger.Logger) error {
	if err := tracing.SetupTracer(cfg.Tracing.Endpoint, cfg.Tracing.ServiceName, ""); err != nil {
		return fmt.Errorf("setup tracer: %w", err)
	}

	msgProcessor, cleanup, err := makeProcessor(
		cfg.Rabbit,
		cfg.Postgres,
		cfg.Profanity,
		cfg.CommandRedis,
		cfg.BanRedis,
		cfg.RankingsRedis,
		cfg.Updates.Token,
	)
	if err != nil {
		return fmt.Errorf("init processor: %w", err)
	}

	defer cleanup()

	tgBot, err := bot.New(cfg.Updates.Token, msgProcessor, l.WithField("bot_name", "profanity_bot"))
	if err != nil {
		return fmt.Errorf("init bot: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()

		l.Info("bot running...")
		tgBot.ProcessUpdates(cfg.Updates.UpdateTimeoutSeconds)
	}()

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)

loop:
	for {
		select {
		case <-terminate:
			signal.Stop(terminate)
			l.Info("stopping bot...")
			tgBot.Stop()
		case <-ctx.Done():
			break loop
		}
	}

	l.Info("bot stopped...")

	return nil
}

func makeProcessor(
	rabbitCfg config.RabbitMQProducer,
	postgresCfg config.Postgres,
	profanityCfg config.Profanity,
	commandRedisCfg config.CommandRedis,
	banCfg config.BanRedis,
	rankingsCfg config.RankingsRedis,
	botToken string,
) (bot.MessageProcessor, func(), error) {
	if rabbitCfg.URL != "" {
		msgPublisher, cleanup, err := makeRabbitPublisher(rabbitCfg)
		if err != nil {
			return nil, nil, fmt.Errorf("rabbit publisher: %w", err)
		}

		return msgPublisher, cleanup, nil
	}

	msgProcessor, cleanup, err := infra.MakeMsgProcessor(
		botToken,
		postgresCfg,
		profanityCfg,
		commandRedisCfg,
		banCfg,
		rankingsCfg,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("msg processor: %w", err)
	}

	return msgProcessor, cleanup, nil
}

func makeRabbitPublisher(rabbitCfg config.RabbitMQProducer) (bot.MessageProcessor, func(), error) {
	ch, cleanup, err := infra.MakeRabbitAMQPChannel(rabbitCfg.URL)
	if err != nil {
		return nil, nil, fmt.Errorf("rabbit channel: %w", err)
	}

	msgPublisher, err := publisher.New(tracing.WrapChannel(ch), rabbitCfg.MsgQueue)
	if err != nil {
		return nil, nil, fmt.Errorf("rabbit publisher: %w", err)
	}

	return msgPublisher, cleanup, nil
}
