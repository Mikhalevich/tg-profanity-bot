package app

import (
	"flag"
	"fmt"
	"os"

	"github.com/jinzhu/configor"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"

	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/msgsender"
	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/profanity"
	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/profanity/matcher"
	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/profanity/replacer"
	"github.com/Mikhalevich/tg-profanity-bot/internal/bot"
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

func MakeMsgProcessor(profanityCfg config.Profanity, botToken string) (bot.MessageProcessor, error) {
	replacer, err := MakeProfanityReplacer(profanityCfg)
	if err != nil {
		return nil, fmt.Errorf("make profanity replacer: %w", err)
	}

	msgSender, err := msgsender.New(botToken)
	if err != nil {
		return nil, fmt.Errorf("make msg sender: %w", err)
	}

	return processor.New(replacer, msgSender), nil
}

// MakeRabbitAMQPChannel make rabbitmq channel and returns channel itself, clearing func and error.
func MakeRabbitAMQPChannel(url string) (*amqp.Channel, func(), error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, fmt.Errorf("ampq dial: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("get channel from connection: %w", err)
	}

	return ch, func() {
		ch.Close()
		conn.Close()
	}, nil
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
