package bot

import (
	"fmt"
	"sync"
	"unicode/utf8"

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
	mangledMsgText := b.replacer.ReplaceMessage(msg.Text)

	if mangledMsgText != msg.Text {
		return b.editMessage(msg, mangledMsgText)
	}

	return nil
}

func (b *bot) editMessage(originMsg *tgbotapi.Message, msgText string) error {
	deletedMsg := tgbotapi.NewDeleteMessage(originMsg.Chat.ID, originMsg.MessageID)
	//nolint:errcheck
	// disabled due to api delete error
	b.api.Send(deletedMsg)

	if _, err := b.api.Send(newMessage(originMsg, msgText)); err != nil {
		return fmt.Errorf("send new: %w", err)
	}

	return nil
}

func newMessage(originMsg *tgbotapi.Message, msgText string) *tgbotapi.MessageConfig {
	formattedMsgText, msgEntities := formatMessage(msgText, originMsg.From)

	newMsg := tgbotapi.NewMessage(originMsg.Chat.ID, formattedMsgText)
	newMsg.Entities = msgEntities

	if originMsg.ReplyToMessage != nil {
		newMsg.ReplyToMessageID = originMsg.ReplyToMessage.MessageID
	}

	return &newMsg
}

func formatMessage(msg string, fromUser *tgbotapi.User) (string, []tgbotapi.MessageEntity) {
	var (
		editedHeader       = "Edited by profanity bot\n"
		editedHeaderOffset = 0
		editedHeaderLen    = utf8.RuneCountInString(editedHeader)
		senderHeader       = "Sender: "
		senderHeaderOffset = editedHeaderOffset + editedHeaderLen
		senderHeaderLen    = utf8.RuneCountInString(senderHeader)
		userName           = extractUserName(fromUser)
		userNameOffset     = senderHeaderOffset + senderHeaderLen
		userNameLen        = utf8.RuneCountInString(userName)
	)

	return fmt.Sprintf("%s%s%s\n%s", editedHeader, senderHeader, fromUser, msg),
		[]tgbotapi.MessageEntity{
			{
				Type:   "bold",
				Offset: editedHeaderOffset,
				Length: editedHeaderLen,
			},
			{
				Type:   "bold",
				Offset: senderHeaderOffset,
				Length: senderHeaderLen,
			},
			{
				Type:   "text_mention",
				Offset: userNameOffset,
				Length: userNameLen,
				User:   fromUser,
			},
		}
}

func extractUserName(from *tgbotapi.User) string {
	if from.UserName != "" {
		return from.UserName
	}

	return fmt.Sprintf("%s %s", from.FirstName, from.LastName)
}

func (b *bot) Stop() {
	b.api.StopReceivingUpdates()
}
