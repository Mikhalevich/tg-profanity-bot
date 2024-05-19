package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/msgsender"
	"github.com/Mikhalevich/tg-profanity-bot/internal/app"
	"github.com/Mikhalevich/tg-profanity-bot/internal/bot"
	"github.com/Mikhalevich/tg-profanity-bot/internal/config"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor"
)

func main() {
	var cfg config.Config
	if err := app.LoadConfig(&cfg); err != nil {
		logrus.WithError(err).Error("failed to load config")
		return
	}

	logger, err := app.SetupLogger(cfg.LogLevel)
	if err != nil {
		logger.WithError(err).Error("failed to setup logger")
		return
	}

	replacer, err := app.MakeProfanityReplacer(cfg.Profanity)
	if err != nil {
		logger.WithError(err).Error("failed to init replacer")
		return
	}

	msgSender, err := msgsender.New(cfg.Bot.Token)
	if err != nil {
		logger.WithError(err).Error("failed to init msg sender")
		return
	}

	msgProcessor := processor.New(replacer, msgSender)

	tgBot, err := bot.New(cfg.Bot.Token, msgProcessor, logger.WithField("bot_name", "profanity_bot"))
	if err != nil {
		logger.WithError(err).Error("configure bot")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()

		logger.Info("bot running...")
		tgBot.ProcessUpdates(cfg.Bot.UpdateTimeoutSeconds)
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
