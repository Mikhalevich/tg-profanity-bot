package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"

	"github.com/Mikhalevich/tg-profanity-bot/internal/app/tracing"
	"github.com/Mikhalevich/tg-profanity-bot/internal/messagequeue/rabbit/internal/contract"
)

type MessageProcessor interface {
	ProcessMessage(ctx context.Context, msg *tgbotapi.Message) error
	ProcessCallbackQuery(ctx context.Context, query *tgbotapi.CallbackQuery) error
}

type consumer struct {
	ch         *amqp.Channel
	queueName  string
	logger     *logrus.Entry
	workerChan chan amqp.Delivery
	wg         sync.WaitGroup
}

func New(ch *amqp.Channel, queueName string, logger *logrus.Entry) (*consumer, error) {
	_, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("declare queue: %w", err)
	}

	return &consumer{
		ch:         ch,
		queueName:  queueName,
		logger:     logger,
		workerChan: make(chan amqp.Delivery),
	}, nil
}

func (c *consumer) Consume(ctx context.Context, workersCount int, processor MessageProcessor) error {
	updates, err := c.ch.ConsumeWithContext(
		ctx,
		c.queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	c.wg.Add(workersCount)
	c.runWorkers(ctx, workersCount, processor)

	for u := range updates {
		c.workerChan <- u
	}

	close(c.workerChan)

	c.logger.Debug("stopping workers")

	c.wg.Wait()

	c.logger.Debug("all workers are stopped")

	return nil
}

func (c *consumer) runWorkers(ctx context.Context, count int, processor MessageProcessor) {
	for range count {
		go func() {
			defer c.wg.Done()

			for d := range c.workerChan {
				if err := c.processUpdate(ctx, d, processor); err != nil {
					c.logger.WithError(err).Error("process update")
					continue
				}
			}
		}()
	}
}

func (c *consumer) processUpdate(ctx context.Context, d amqp.Delivery, processor MessageProcessor) error {
	ctx = tracing.ExtractContextFromHeaders(ctx, d.Headers)

	c.logger.WithField("message_type", d.Type).Debug("received rabbit update")

	switch contract.MessageType(d.Type) {
	case contract.Message:
		return c.processMessage(ctx, d.Body, processor)

	case contract.CallbackQuery:
		return c.processCallbackQuery(ctx, d.Body, processor)
	}

	return fmt.Errorf("unsupported type: %s", d.Type)
}

func (c *consumer) processMessage(ctx context.Context, body []byte, processor MessageProcessor) error {
	var msg *tgbotapi.Message
	if err := json.Unmarshal(body, &msg); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	if err := processor.ProcessMessage(ctx, msg); err != nil {
		return fmt.Errorf("process message: %w", err)
	}

	return nil
}

func (c *consumer) processCallbackQuery(ctx context.Context, body []byte, processor MessageProcessor) error {
	var query *tgbotapi.CallbackQuery
	if err := json.Unmarshal(body, &query); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	if err := processor.ProcessCallbackQuery(ctx, query); err != nil {
		return fmt.Errorf("process callback query: %w", err)
	}

	return nil
}
