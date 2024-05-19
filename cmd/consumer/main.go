package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jinzhu/configor"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"

	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/msgsender"
	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/profanity"
	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/profanity/matcher"
	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/profanity/replacer"
	"github.com/Mikhalevich/tg-profanity-bot/internal/config"
	"github.com/Mikhalevich/tg-profanity-bot/internal/messagequeue/rabbit/consumer"
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

	msgSender, err := msgsender.New(cfg.BotToken)
	if err != nil {
		logger.WithError(err).Error("failed to init msg sender")
		return
	}

	msgProcessor := processor.New(replacer, msgSender)

	conn, err := amqp.Dial(cfg.Rabbit.URL)
	if err != nil {
		logger.WithError(err).Error("failed to init rabbit mq connection")
		return
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logger.WithError(err).Error("failed to init rabbit mq channel")
		return
	}

	defer ch.Close()

	ctx, cancel := context.WithCancel(context.Background())

	c, err := consumer.New(ch, cfg.Rabbit.MsgQueue, logger.WithField("bot_name", "bot_msg_worker"))
	if err != nil {
		logger.WithError(err).Error("failed to init message consumer")
		return
	}

	go func() {
		defer cancel()

		terminate := make(chan os.Signal, 1)
		signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)

		<-terminate
	}()

	logger.Info("consumer running...")

	if err := c.Consume(ctx, 10, msgProcessor); err != nil {
		logger.WithError(err).Error("consume messages")
		return
	}

	logger.Info("consumer stopped...")
}

func loadConfig() (*config.Consumer, error) {
	configFile := flag.String("config", "config/config-consumer.yaml", "consumer worker config file")
	flag.Parse()

	var cfg config.Consumer
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
