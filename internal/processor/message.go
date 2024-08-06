package processor

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Mikhalevich/tg-profanity-bot/internal/app/tracing"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/cmd"
)

func (p *processor) ProcessMessage(ctx context.Context, msg *tgbotapi.Message) error {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	var (
		chatID = strconv.FormatInt(msg.Chat.ID, 10)
		userID = strconv.FormatInt(msg.From.ID, 10)
	)

	if p.banProcessor.IsBanned(ctx, makeBanID(chatID, userID)) {
		if err := p.processBan(ctx, msg); err != nil {
			return fmt.Errorf("process ban: %w", err)
		}

		return nil
	}

	isProcessed, err := p.tryProcessCommand(ctx, chatID, msg)
	if err != nil {
		return fmt.Errorf("process command: %w", err)
	}

	if isProcessed {
		return nil
	}

	if err := p.processReplace(ctx, chatID, userID, msg); err != nil {
		return fmt.Errorf("process replace: %w", err)
	}

	return nil
}

func makeBanID(chatID, userID string) string {
	return fmt.Sprintf("%s:%s", chatID, userID)
}

func (p *processor) processReplace(ctx context.Context, chatID, userID string, msg *tgbotapi.Message) error {
	mangledMsg, err := p.replacer.Replace(ctx, chatID, msg.Text)
	if err != nil {
		return fmt.Errorf("replace msg: %w", err)
	}

	if mangledMsg == msg.Text {
		return nil
	}

	if err := p.msgSender.Edit(ctx, msg, mangledMsg, p.viewOriginMsgButton(ctx, msg.Text)...); err != nil {
		return fmt.Errorf("msg edit: %w", err)
	}

	if err := p.banProcessor.AddViolation(ctx, makeBanID(chatID, userID)); err != nil {
		return fmt.Errorf("add violation: %w", err)
	}

	return nil
}

func (p *processor) viewOriginMsgButton(ctx context.Context, msg string) []tgbotapi.InlineKeyboardButton {
	return p.makeButton(ctx, "view origin msg", Command{
		CMD:     cmd.ViewOrginMsg.String(),
		Payload: []byte(msg),
	})
}

func (p *processor) viewBannedMsgButton(ctx context.Context, msg string) []tgbotapi.InlineKeyboardButton {
	return p.makeButton(ctx, "view origin msg", Command{
		CMD:     cmd.ViewBannedMsg.String(),
		Payload: []byte(msg),
	})
}

func (p *processor) processBan(ctx context.Context, msg *tgbotapi.Message) error {
	if err := p.msgSender.Edit(
		ctx,
		msg,
		"user is banned",
		p.viewBannedMsgButton(ctx, msg.Text)...,
	); err != nil {
		return fmt.Errorf("msg edit: %w", err)
	}

	return nil
}
