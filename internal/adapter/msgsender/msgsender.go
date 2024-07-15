package msgsender

import (
	"context"
	"fmt"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Mikhalevich/tg-profanity-bot/internal/app/tracing"
)

type msgsender struct {
	api *tgbotapi.BotAPI
}

func New(token string) (*msgsender, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("create bot api: %w", err)
	}

	return &msgsender{
		api: api,
	}, nil
}

func (s *msgsender) Reply(ctx context.Context, originMsg *tgbotapi.Message, msg string) error {
	_, span := tracing.StartSpan(ctx)
	defer span.End()

	newMsg := tgbotapi.NewMessage(originMsg.Chat.ID, msg)
	newMsg.ReplyToMessageID = originMsg.MessageID

	if _, err := s.api.Send(newMsg); err != nil {
		return fmt.Errorf("send reply: %w", err)
	}

	return nil
}

func (s *msgsender) Edit(ctx context.Context, originMsg *tgbotapi.Message, msg string) error {
	_, span := tracing.StartSpan(ctx)
	defer span.End()

	deletedMsg := tgbotapi.NewDeleteMessage(originMsg.Chat.ID, originMsg.MessageID)
	//nolint:errcheck
	// disabled due to api delete error
	s.api.Send(deletedMsg)

	if _, err := s.api.Send(newEditedMessage(originMsg, msg)); err != nil {
		return fmt.Errorf("send new: %w", err)
	}

	return nil
}

func newEditedMessage(originMsg *tgbotapi.Message, msgText string) *tgbotapi.MessageConfig {
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
