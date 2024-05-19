package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"

	"github.com/Mikhalevich/tg-profanity-bot/internal/adapter/msgsender"
	"github.com/Mikhalevich/tg-profanity-bot/internal/app"
	"github.com/Mikhalevich/tg-profanity-bot/internal/config"
	"github.com/Mikhalevich/tg-profanity-bot/internal/messagequeue/rabbit/consumer"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor"
)

func main() {
	var cfg config.Consumer
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
