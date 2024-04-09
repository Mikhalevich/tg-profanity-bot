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

	"github.com/Mikhalevich/tg-profanity-bot/internal/bot"
	"github.com/Mikhalevich/tg-profanity-bot/internal/profanity"
)

type config struct {
	BotToken string `yaml:"bot_token" required:"true"`
}

func main() {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	cfg, err := loadConfig()
	if err != nil {
		logger.WithError(err).Error("failed to load config")
		return
	}

	tgBot, err := bot.New(cfg.BotToken, isDebugEnabled(), profanity.New(), logger.WithField("bot_name", "profanity_bot"))
	if err != nil {
		logger.WithError(err).Error("configure bot")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()

		logger.Info("bot running...")
		tgBot.ProcessUpdates(5)
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

func isDebugEnabled() bool {
	if debug := os.Getenv("FB_DEBUG"); debug != "" {
		return true
	}

	return false
}

func loadConfig() (*config, error) {
	configFile := flag.String("config", "config/config.yaml", "telegram bot config file")
	flag.Parse()

	var cfg config
	if err := configor.Load(&cfg, *configFile); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	return &cfg, nil
}
