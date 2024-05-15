package publisher

import (
	"context"
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	amqp "github.com/rabbitmq/amqp091-go"
)

type publisher struct {
	ch        *amqp.Channel
	queueName string
}

func New(ch *amqp.Channel, queueName string) (*publisher, error) {
	_, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("declare queue: %w", err)
	}

	return &publisher{
		ch:        ch,
		queueName: queueName,
	}, nil
}

func (p *publisher) ProcessMessage(msg *tgbotapi.Message) error {
	if err := p.publishJSON(context.Background(), msg); err != nil {
		return fmt.Errorf("publish json: %w", err)
	}

	return nil
}

func (p *publisher) publishJSON(ctx context.Context, obj any) error {
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
		},
	); err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	return nil
}
