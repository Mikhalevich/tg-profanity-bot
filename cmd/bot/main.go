package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"

	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/msgsender"
	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/profanity"
	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/profanity/matcher"
	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/profanity/replacer"
	"github.com/Mikhalevich/tg-profanity-bot/internal/bot"
	"github.com/Mikhalevich/tg-profanity-bot/internal/config"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		logrus.WithError(err).Error("failed to load config")
		return
	}

	logger, err := setupLogger(cfg.LogLevel)
	if err != nil {
		logger.WithError(err).Error("failed to setup logger")
		return
	}

	replacer, err := makeProfanityReplacer(cfg.Profanity)
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

func loadConfig() (*config.Config, error) {
	configFile := flag.String("config", "config/config.yaml", "telegram bot config file")
	flag.Parse()

	var cfg config.Config
	if err := configor.Load(&cfg, *configFile); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	return &cfg, nil
}

func setupLogger(lvl string) (*logrus.Logger, error) {
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		return nil, fmt.Errorf("parse log level: %w", err)
	}

	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return logger, nil
}

func makeProfanityReplacer(cfg config.Profanity) (processor.TextReplacer, error) {
	words, err := config.BadWords()
	if err != nil {
		return nil, fmt.Errorf("get bad words: %w", err)
	}

	if cfg.Dynamic != "" {
		return profanity.New(matcher.NewAhocorasick(words), replacer.NewDynamic(cfg.Dynamic)), nil
	}

	return profanity.New(matcher.NewAhocorasick(words), replacer.NewStatic(cfg.Static)), nil
}
