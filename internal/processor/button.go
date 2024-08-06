package processor

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/cmd"
)

func (p *processor) makeButton(ctx context.Context, caption string, command Command) *tgbotapi.InlineKeyboardButton {
	id := uuid.NewString()

	if err := p.commandStorage.Set(ctx, id, command); err != nil {
		// skip error
		return nil
	}

	button := tgbotapi.NewInlineKeyboardButtonData(caption, id)

	return &button
}

func buttonRow(buttons ...*tgbotapi.InlineKeyboardButton) []tgbotapi.InlineKeyboardButton {
	row := make([]tgbotapi.InlineKeyboardButton, 0, len(buttons))

	for _, b := range buttons {
		if b != nil {
			row = append(row, *b)
		}
	}

	if len(row) == 0 {
		return nil
	}

	return tgbotapi.NewInlineKeyboardRow(row...)
}

func (p *processor) revertButton(ctx context.Context, command cmd.CMD, word string) *tgbotapi.InlineKeyboardButton {
	return p.makeButton(ctx, "revert", Command{
		CMD:     command.String(),
		Payload: []byte(word),
	})
}

func (p *processor) viewOriginMsgButton(ctx context.Context, msg string) *tgbotapi.InlineKeyboardButton {
	return p.makeButton(ctx, "view origin msg", Command{
		CMD:     cmd.ViewOrginMsg.String(),
		Payload: []byte(msg),
	})
}

func (p *processor) viewBannedMsgButton(ctx context.Context, msg string) *tgbotapi.InlineKeyboardButton {
	return p.makeButton(ctx, "view origin msg", Command{
		CMD:     cmd.ViewBannedMsg.String(),
		Payload: []byte(msg),
	})
}

func (p *processor) unbanButton(ctx context.Context, userID string) *tgbotapi.InlineKeyboardButton {
	return p.makeButton(ctx, "unban user", Command{
		CMD:     cmd.Unban.String(),
		Payload: []byte(userID),
	})
}
