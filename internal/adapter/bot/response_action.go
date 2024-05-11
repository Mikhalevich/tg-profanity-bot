package bot

import (
	"fmt"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type responseAction struct {
	originMsg *tgbotapi.Message
	api       *tgbotapi.BotAPI
}

func newResponseAction(originMsg *tgbotapi.Message, api *tgbotapi.BotAPI) *responseAction {
	return &responseAction{
		originMsg: originMsg,
		api:       api,
	}
}

func (ra *responseAction) Send(msg string) error {
	return nil
}

func (ra *responseAction) Edit(msg string) error {
	deletedMsg := tgbotapi.NewDeleteMessage(ra.originMsg.Chat.ID, ra.originMsg.MessageID)
	//nolint:errcheck
	// disabled due to api delete error
	ra.api.Send(deletedMsg)

	if _, err := ra.api.Send(newMessage(ra.originMsg, msg)); err != nil {
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
