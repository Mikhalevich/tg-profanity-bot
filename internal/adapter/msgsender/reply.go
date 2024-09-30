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
	options ...port.Option,
) error {
	_, span := tracing.StartSpan(ctx)
	defer span.End()

	var (
		opts     = makeOptions(options)
		entities = convertToMessageEntities(opts.Format)
	)

	newMsg := tgbotapi.NewMessage(originMsgInfo.ChatID.Int64(), msg)
	newMsg.ReplyToMessageID = originMsgInfo.MessageID

	if len(entities) > 0 {
		newMsg.Entities = entities
	}

	if len(opts.Buttons) > 0 {
		newMsg.ReplyMarkup = buttonsMarkup(opts.Buttons)
	}

	if _, err := s.api.Send(newMsg); err != nil {
		return fmt.Errorf("send reply: %w", err)
	}

	return nil
}

func convertToMessageEntities(format []port.Format) []tgbotapi.MessageEntity {
	if len(format) == 0 {
		return nil
	}

	entities := make([]tgbotapi.MessageEntity, 0, len(format))

	for _, f := range format {
		entities = append(entities, tgbotapi.MessageEntity{
			Type:   f.Type.String(),
			Offset: f.Offset,
			Length: f.Length,
			User:   f.User,
		})
	}

	return entities
}
