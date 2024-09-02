package bot

import (
	"context"
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Mikhalevich/tg-profanity-bot/internal/app/logger"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

type MessageProcessor interface {
	ProcessMessage(ctx context.Context, info port.MessageInfo) error
	ProcessCallbackQuery(ctx context.Context, query port.CallbackQuery) error
}

type bot struct {
	api       *tgbotapi.BotAPI
	processor MessageProcessor
	logger    logger.Logger
}

func New(
	token string,
	processor MessageProcessor,
	logger logger.Logger,
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
		wg.Add(1)

		go func(update *tgbotapi.Update) {
			defer wg.Done()

			ctx := logger.WithLogger(context.Background(), b.logger)
			if err := b.processUpdate(ctx, update); err != nil {
				b.logger.WithError(err).Error("process update")
			}
		}(&update)
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

func toMessageInfo(msg *tgbotapi.Message) port.MessageInfo {
	return port.MessageInfo{
		MessageID:        msg.MessageID,
		ChatID:           port.NewID(msg.Chat.ID),
		UserID:           port.NewID(msg.From.ID),
		UserFrom:         msg.From,
		Text:             msg.Text,
		ReplyToMessageID: extractReplyToMessageID(msg),
	}
}

func extractReplyToMessageID(msg *tgbotapi.Message) int {
	if msg.ReplyToMessage != nil {
		return msg.ReplyToMessage.MessageID
	}

	return 0
}

func toCallbackQuery(query *tgbotapi.CallbackQuery) port.CallbackQuery {
	return port.CallbackQuery{
		MessageInfo: toMessageInfo(query.Message),
		Data:        query.Data,
	}
}

func (b *bot) processUpdate(ctx context.Context, update *tgbotapi.Update) error {
	msg := extractTextMessage(update)
	if msg != nil {
		b.logger.WithFields(logger.Fields{
			"chat_id": msg.Chat.ID,
			"message": msg.Text,
		}).Debug("incoming message")

		if err := b.processor.ProcessMessage(ctx, toMessageInfo(msg)); err != nil {
			return fmt.Errorf("process message: %w", err)
		}

		return nil
	}

	if update.CallbackQuery != nil {
		if err := b.processor.ProcessCallbackQuery(ctx, toCallbackQuery(update.CallbackQuery)); err != nil {
			return fmt.Errorf("process callback query: %w", err)
		}

		return nil
	}

	return nil
}

func (b *bot) Stop() {
	b.api.StopReceivingUpdates()
}
