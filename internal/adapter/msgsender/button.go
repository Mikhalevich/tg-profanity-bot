package msgsender

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func buttonsMarkup(buttons []*port.Button) tgbotapi.InlineKeyboardMarkup {
	row := make([]tgbotapi.InlineKeyboardButton, 0, len(buttons))

	for _, b := range buttons {
		if b != nil {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(b.Caption, b.Data))
		}
	}

	if len(row) == 0 {
		return tgbotapi.NewInlineKeyboardMarkup()
	}

	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(row...))
}
