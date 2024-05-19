package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/Mikhalevich/tg-profanity-bot/internal/app"
	"github.com/Mikhalevich/tg-profanity-bot/internal/bot"
	"github.com/Mikhalevich/tg-profanity-bot/internal/config"
)

func main() {
	var cfg config.Bot
	if err := app.LoadConfig(&cfg); err != nil {
		logrus.WithError(err).Error("failed to load config")
		return
	}

	logger, err := app.SetupLogger(cfg.LogLevel)
	if err != nil {
		logger.WithError(err).Error("failed to setup logger")
		return
	}

	msgProcessor, err := app.MakeMsgProcessor(cfg.Profanity, cfg.Updates.Token)
	if err != nil {
		logger.WithError(err).Error("failed to init msg processor")
		return
	}

	tgBot, err := bot.New(cfg.Updates.Token, msgProcessor, logger.WithField("bot_name", "profanity_bot"))
	if err != nil {
		logger.WithError(err).Error("configure bot")
		return
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
}
