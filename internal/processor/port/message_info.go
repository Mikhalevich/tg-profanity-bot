package port

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ID struct {
	IdNum int64
	IdStr string
}

func NewID(id int64) ID {
	return ID{
		IdNum: id,
		IdStr: strconv.FormatInt(id, 10),
	}
}

func (id ID) Int64() int64 {
	return id.IdNum
}

func (id ID) String() string {
	return id.IdStr
}

type MessageInfo struct {
	MessageID        int
	ChatID           ID
	UserID           ID
	UserFrom         *tgbotapi.User
	Text             string
	ReplyToMessageID int
}

type CallbackQuery struct {
	MessageInfo
	Data string
}
