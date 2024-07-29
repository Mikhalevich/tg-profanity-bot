package processor

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Mikhalevich/tg-profanity-bot/internal/app/tracing"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/button"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/cmd"
)

func (p *processor) ProcessMessage(ctx context.Context, msg *tgbotapi.Message) error {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	var (
		chatID = strconv.FormatInt(msg.Chat.ID, 10)
	)

	isProcessed, err := p.tryProcessCommand(ctx, chatID, msg)
	if err != nil {
		return fmt.Errorf("process command: %w", err)
	}

	if isProcessed {
		return nil
	}

	if err := p.processReplace(ctx, chatID, msg); err != nil {
		return fmt.Errorf("process replace: %w", err)
	}

	return nil
}

func (p *processor) processReplace(ctx context.Context, chatID string, msg *tgbotapi.Message) error {
	mangledMsg, err := p.replacer.Replace(ctx, chatID, msg.Text)
	if err != nil {
		return fmt.Errorf("replace msg: %w", err)
	}

	if mangledMsg == msg.Text {
		return nil
	}

	if err := p.msgSender.Edit(ctx, msg, mangledMsg, viewOriginMsgButton(msg.Text)); err != nil {
		return fmt.Errorf("msg edit: %w", err)
	}

	return nil
}

func viewOriginMsgButton(msg string) []tgbotapi.InlineKeyboardButton {
	buttonInfo, err := button.ButtonCMDInfo{
		CMD:     cmd.ViewOrginMsg.String(),
		Payload: []byte(msg),
	}.ToBase64()

	if err != nil {
		return nil
	}

	return tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("view origin msg", buttonInfo),
	)
}
