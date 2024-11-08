package publisher

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/Mikhalevich/tg-profanity-bot/internal/app/messagequeue/rabbit/internal/contract"
	"github.com/Mikhalevich/tg-profanity-bot/internal/infra/tracing"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
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

func (p *publisher) ProcessMessage(ctx context.Context, info port.MessageInfo) error {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	if err := p.publishGOB(ctx, info, contract.Message); err != nil {
		return fmt.Errorf("publish gob: %w", err)
	}

	return nil
}

func (p *publisher) ProcessCallbackQuery(ctx context.Context, query port.CallbackQuery) error {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	if err := p.publishGOB(ctx, query, contract.CallbackQuery); err != nil {
		return fmt.Errorf("publish gob: %w", err)
	}

	return nil
}

func (p *publisher) publishGOB(ctx context.Context, obj any, messageType contract.MessageType) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(obj); err != nil {
		return fmt.Errorf("gob encode: %w", err)
	}

	if err := p.ch.PublishWithContext(
		ctx,
		"",
		p.queueName,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "gob",
			Body:         buf.Bytes(),
			Type:         messageType.String(),
		},
	); err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	return nil
}
