package publisher

import (
	"context"
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/Mikhalevich/tg-profanity-bot/internal/app/tracing"
	"github.com/Mikhalevich/tg-profanity-bot/internal/messagequeue/rabbit/internal/contract"
)

type ChannelPublisher interface {
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	PublishWithContext(
		ctx context.Context,
		exchange string,
		key string,
		mandatory bool,
		immediate bool,
		msg amqp.Publishing,
	) error
}

type publisher struct {
	ch        ChannelPublisher
	queueName string
}

func New(ch ChannelPublisher, queueName string) (*publisher, error) {
	_, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("declare queue: %w", err)
	}

	return &publisher{
		ch:        ch,
		queueName: queueName,
	}, nil
}

func (p *publisher) ProcessMessage(ctx context.Context, msg *tgbotapi.Message) error {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	if err := p.publishJSON(ctx, msg, contract.Message); err != nil {
		return fmt.Errorf("publish json: %w", err)
	}

	return nil
}

func (p *publisher) ProcessCallbackQuery(ctx context.Context, query *tgbotapi.CallbackQuery) error {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	if err := p.publishJSON(ctx, query, contract.CallbackQuery); err != nil {
		return fmt.Errorf("publish json: %w", err)
	}

	return nil
}

func (p *publisher) publishJSON(ctx context.Context, obj any, messageType contract.MessageType) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	if err := p.ch.PublishWithContext(
		ctx,
		"",
		p.queueName,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         b,
			Type:         messageType.String(),
		},
	); err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	return nil
}
