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

func New(api *tgbotapi.BotAPI) *msgsender {
	return &msgsender{
		api: api,
	}
}

func (s *msgsender) Edit(
	ctx context.Context,
	originMsg *tgbotapi.Message,
	msg string,
	buttons ...tgbotapi.InlineKeyboardButton,
) error {
	_, span := tracing.StartSpan(ctx)
	defer span.End()

	deletedMsg := tgbotapi.NewDeleteMessage(originMsg.Chat.ID, originMsg.MessageID)
	//nolint:errcheck
	// disabled due to api delete error
	s.api.Send(deletedMsg)

	if _, err := s.api.Send(newEditedMessage(originMsg, msg, buttons)); err != nil {
		return fmt.Errorf("send new: %w", err)
	}

	return nil
}

func newEditedMessage(
	originMsg *tgbotapi.Message,
	msgText string,
	buttons []tgbotapi.InlineKeyboardButton,
) *tgbotapi.MessageConfig {
	formattedMsgText, msgEntities := formatMessage(msgText, originMsg.From)

	newMsg := tgbotapi.NewMessage(originMsg.Chat.ID, formattedMsgText)
	newMsg.Entities = msgEntities

	if len(buttons) > 0 {
		newMsg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons)
	}

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
