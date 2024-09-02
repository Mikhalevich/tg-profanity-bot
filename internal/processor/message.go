package processor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-profanity-bot/internal/app/tracing"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (p *processor) ProcessMessage(ctx context.Context, info port.MessageInfo) error {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	if p.banProcessor.IsBanned(ctx, makeBanID(info.ChatID.String(), info.UserID.String())) {
		if err := p.processBan(ctx, info.UserID.String(), info); err != nil {
			return fmt.Errorf("process ban: %w", err)
		}

		return nil
	}

	isProcessed, err := p.tryProcessCommand(ctx, info)
	if err != nil {
		return fmt.Errorf("process command: %w", err)
	}

	if isProcessed {
		return nil
	}

	if err := p.processReplace(ctx, info); err != nil {
		return fmt.Errorf("process replace: %w", err)
	}

	return nil
}

func makeBanID(chatID, userID string) string {
	return fmt.Sprintf("%s:%s", chatID, userID)
}

func (p *processor) processReplace(ctx context.Context, info port.MessageInfo) error {
	mangledMsg, err := p.mangler.Mangle(ctx, info.ChatID.String(), info.Text)
	if err != nil {
		return fmt.Errorf("replace msg: %w", err)
	}

	if mangledMsg == info.Text {
		return nil
	}

	isBanned, err := p.banProcessor.AddViolation(ctx, makeBanID(info.ChatID.String(), info.UserID.String()))
	if err != nil {
		return fmt.Errorf("add violation: %w", err)
	}

	if isBanned {
		if err := p.processBan(ctx, info.UserID.String(), info); err != nil {
			return fmt.Errorf("process ban: %w", err)
		}

		return nil
	}

	if err := p.msgSender.Edit(
		ctx,
		info,
		mangledMsg,
		p.viewOriginMsgButton(ctx, info.Text),
	); err != nil {
		return fmt.Errorf("msg edit: %w", err)
	}

	return nil
}

func (p *processor) processBan(ctx context.Context, userID string, info port.MessageInfo) error {
	if err := p.msgSender.Edit(
		ctx,
		info,
		"user is banned",
		p.viewBannedMsgButton(ctx, info.Text),
		p.unbanButton(ctx, userID),
	); err != nil {
		return fmt.Errorf("msg edit: %w", err)
	}

	return nil
}
