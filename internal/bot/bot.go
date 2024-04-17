package bot

import (
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type MessageReplacer interface {
	ReplaceMessage(msg string) string
}

type bot struct {
	api      *tgbotapi.BotAPI
	replacer MessageReplacer
	logger   *logrus.Entry
}

func New(
	token string,
	replacer MessageReplacer,
	logger *logrus.Entry,
) (*bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("create bot api: %w", err)
	}

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
			continue
		}

		b.logger.WithFields(logrus.Fields{
			"chat_id": update.Message.Chat.ID,
			"message": update.Message.Text,
		}).Debug("incoming message")

		wg.Add(1)

		go func(u tgbotapi.Update) {
			defer wg.Done()

			if err := b.processMessage(u.Message); err != nil {
				b.logger.WithError(err).Error("process message")
			}
		}(update)
	}

	wg.Wait()
}

func (b *bot) processMessage(msg *tgbotapi.Message) error {
	mangledMsg := b.replacer.ReplaceMessage(msg.Text)

	if mangledMsg != msg.Text {
		formattedMsg := formatMessage(mangledMsg, userName(msg.From))
		return b.editMessage(msg, formattedMsg)
	}

	return nil
}

func userName(from *tgbotapi.User) string {
	if from.UserName != "" {
		return from.UserName
	}

	return fmt.Sprintf("%s %s", from.FirstName, from.LastName)
}

func formatMessage(msg string, fromUser string) string {
	return fmt.Sprintf("Edited by profanity bot\nSender: %s\n\n%s", fromUser, msg)
}

func (b *bot) editMessage(originMsg *tgbotapi.Message, text string) error {
	deletedMsg := tgbotapi.NewDeleteMessage(originMsg.Chat.ID, originMsg.MessageID)
	//nolint:errcheck
	// disabled due to api delete error
	b.api.Send(deletedMsg)

	newMsg := tgbotapi.NewMessage(originMsg.Chat.ID, text)
	if originMsg.ReplyToMessage != nil {
		newMsg.ReplyToMessageID = originMsg.ReplyToMessage.MessageID
	}

	if _, err := b.api.Send(newMsg); err != nil {
		return fmt.Errorf("send new: %w", err)
	}

	return nil
}

func (b *bot) Stop() {
	b.api.StopReceivingUpdates()
}
