package msgsender

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Mikhalevich/tg-profanity-bot/internal/app/tracing"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (s *msgsender) Reply(
	ctx context.Context,
	originMsgInfo port.MessageInfo,
	msg string,
	buttons ...tgbotapi.InlineKeyboardButton,
) error {
	_, span := tracing.StartSpan(ctx)
	defer span.End()

	newMsg := tgbotapi.NewMessage(originMsgInfo.ChatID.Int64(), msg)
	newMsg.ReplyToMessageID = originMsgInfo.MessageID

	if len(buttons) > 0 {
		newMsg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons)
	}

	if _, err := s.api.Send(newMsg); err != nil {
		return fmt.Errorf("send reply: %w", err)
	}

	return nil
}
