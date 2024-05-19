package app

import (
	"flag"
	"fmt"
	"os"

	"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"

	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/profanity"
	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/profanity/matcher"
	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/profanity/replacer"
	"github.com/Mikhalevich/tg-profanity-bot/internal/config"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor"
)

func MakeProfanityReplacer(cfg config.Profanity) (processor.TextReplacer, error) {
	words, err := config.BadWords()
	if err != nil {
		return nil, fmt.Errorf("get bad words: %w", err)
	}

	if cfg.Dynamic != "" {
		return profanity.New(matcher.NewAhocorasick(words), replacer.NewDynamic(cfg.Dynamic)), nil
	}

	return profanity.New(matcher.NewAhocorasick(words), replacer.NewStatic(cfg.Static)), nil
}

func LoadConfig(cfg any) error {
	configFile := flag.String("config", "config/config.yaml", "consumer worker config file")
	flag.Parse()

	if err := configor.Load(cfg, *configFile); err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	return nil
}

func SetupLogger(lvl string) (*logrus.Logger, error) {
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		return nil, fmt.Errorf("parse log level: %w", err)
	}

	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})

	return logger, nil
}
