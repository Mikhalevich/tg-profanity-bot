package msgsender

import (
	"context"
	"fmt"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Mikhalevich/tg-profanity-bot/internal/app/tracing"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

var _ port.MsgSender = (*msgsender)(nil)

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
	originMsgInfo port.MessageInfo,
	msg string,
	buttons ...*port.Button,
) error {
	_, span := tracing.StartSpan(ctx)
	defer span.End()

	deletedMsg := tgbotapi.NewDeleteMessage(originMsgInfo.ChatID.Int64(), originMsgInfo.MessageID)
	//nolint:errcheck
	// disabled due to api delete error
	s.api.Send(deletedMsg)

	if _, err := s.api.Send(newEditedMessage(originMsgInfo, msg, buttons)); err != nil {
		return fmt.Errorf("send new: %w", err)
	}

	return nil
}

func newEditedMessage(
	originMsgInfo port.MessageInfo,
	msgText string,
	buttons []*port.Button,
) *tgbotapi.MessageConfig {
	formattedMsgText, msgEntities := formatMessage(msgText, originMsgInfo.UserFrom)

	newMsg := tgbotapi.NewMessage(originMsgInfo.ChatID.Int64(), formattedMsgText)
	newMsg.Entities = msgEntities

	if len(buttons) > 0 {
		newMsg.ReplyMarkup = buttonsMarkup(buttons)
	}

	if originMsgInfo.ReplyToMessageID != 0 {
		newMsg.ReplyToMessageID = originMsgInfo.ReplyToMessageID
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
		userName           = fromUser.String()
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
