package bot

import (
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor"
)

type MessageProcessor interface {
	ProcessMessage(msg string, rsp processor.ResponseAction) error
	ProcessCommand(cmd string, args string, rsp processor.ResponseAction) error
}

type bot struct {
	api       *tgbotapi.BotAPI
	processor MessageProcessor
	logger    *logrus.Entry
}

func New(
	token string,
	processor MessageProcessor,
	logger *logrus.Entry,
) (*bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("create bot api: %w", err)
	}

	return &bot{
		api:       api,
		processor: processor,
		logger:    logger,
	}, nil
}

func (b *bot) ProcessUpdates(timeout int) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = timeout

	updates := b.api.GetUpdatesChan(u)

	var wg sync.WaitGroup

	for update := range updates {
		msg := extractTextMessage(&update)
		if msg == nil {
			continue
		}

		b.logger.WithFields(logrus.Fields{
			"chat_id": msg.Chat.ID,
			"message": msg.Text,
		}).Debug("incoming message")

		wg.Add(1)

		go func(msg *tgbotapi.Message) {
			defer wg.Done()

			if err := b.processMessage(msg); err != nil {
				b.logger.WithError(err).Error("process message")
			}
		}(msg)
	}

	wg.Wait()
}

func extractTextMessage(u *tgbotapi.Update) *tgbotapi.Message {
	msg := extractMessage(u)
	if msg != nil && msg.Text != "" {
		return msg
	}

	return nil
}

func extractMessage(u *tgbotapi.Update) *tgbotapi.Message {
	if u.Message != nil {
		return u.Message
	}

	if u.EditedMessage != nil {
		return u.EditedMessage
	}

	return nil
}

func (b *bot) processMessage(msg *tgbotapi.Message) error {
	if err := b.processor.ProcessMessage(msg.Text, newResponseAction(msg, b.api)); err != nil {
		return fmt.Errorf("process message: %w", err)
	}

	return nil
}

func (b *bot) Stop() {
	b.api.StopReceivingUpdates()
}
