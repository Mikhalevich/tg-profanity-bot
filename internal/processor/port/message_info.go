package port

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ID struct {
	IDNum int64
	IDStr string
}

func NewID(id int64) ID {
	return ID{
		IDNum: id,
		IDStr: strconv.FormatInt(id, 10),
	}
}

func NewIDFromString(id string) (ID, error) {
	idNum, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return ID{}, fmt.Errorf("parse int: %w", err)
	}

	return ID{
		IDNum: idNum,
		IDStr: id,
	}, nil
}

func (id ID) Int64() int64 {
	return id.IDNum
}

func (id ID) String() string {
	return id.IDStr
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
