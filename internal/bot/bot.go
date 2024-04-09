package bot

import (
	"context"
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type MessageReplacer interface {
	ReplaceMessage(ctx context.Context, msg string) (string, error)
}

type bot struct {
	api      *tgbotapi.BotAPI
	replacer MessageReplacer
	logger   *logrus.Entry
}

func New(
	token string,
	isDebugEnabled bool,
	replacer MessageReplacer,
	logger *logrus.Entry,
) (*bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("create bot api: %w", err)
	}

	api.Debug = isDebugEnabled

	return &bot{
		api:      api,
		replacer: replacer,
		logger:   logger,
	}, nil
}

func (b *bot) ProcessUpdates(timeout int) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = timeout

	updates := b.api.GetUpdatesChan(u)

	var wg sync.WaitGroup

	for update := range updates {
		if update.Message == nil || update.Message.Text == "" {
			b.logger.Error("invalid message")
			continue
		}

		b.logger.WithFields(logrus.Fields{
			"chat_id": update.Message.Chat.ID,
			"message": update.Message.Text,
		}).Debug("incoming message")

		wg.Add(1)

		go func(u tgbotapi.Update) {
			defer wg.Done()

			if err := b.processMessage(
				context.Background(),
				u.Message.MessageID,
				u.Message.Chat.ID,
				u.Message.Text,
			); err != nil {
				b.logger.WithError(err).Error("process message")
			}
		}(update)
	}

	wg.Wait()
}

func (b *bot) processMessage(ctx context.Context, messageID int, chatID int64, msg string) error {
	output, err := b.replacer.ReplaceMessage(ctx, msg)
	if err != nil {
		return fmt.Errorf("process message: %w", err)
	}

	if output != "" {
		return b.editMessage(messageID, chatID, output)
	}

	return nil
}

func (b *bot) editMessage(messageID int, chatID int64, text string) error {
	msg := tgbotapi.NewEditMessageText(chatID, messageID, text)

	if _, err := b.api.Send(msg); err != nil {
		return fmt.Errorf("send: %w", err)
	}

	return nil
}

func (b *bot) Stop() {
	b.api.StopReceivingUpdates()
}
