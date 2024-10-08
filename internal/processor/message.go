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

	if err := p.updateRankingScore(ctx, info); err != nil {
		return fmt.Errorf("update ranking score: %w", err)
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
		port.WithButton(p.viewOriginMsgButton(ctx, info.Text)),
	); err != nil {
		return fmt.Errorf("msg edit: %w", err)
	}

	return nil
}

func (p *processor) updateRankingScore(ctx context.Context, info port.MessageInfo) error {
	if err := p.rankings.AddScore(
		ctx,
		makeCurrentMonthRankingKey(info.ChatID.String()),
		info.UserID.String(),
	); err != nil {
		return fmt.Errorf("add rankings score: %w", err)
	}

	return nil
}

func (p *processor) processBan(ctx context.Context, userID string, info port.MessageInfo) error {
	if err := p.msgSender.Edit(
		ctx,
		info,
		"user is banned",
		port.WithButton(p.viewBannedMsgButton(ctx, info.Text)),
		port.WithButton(p.unbanButton(ctx, userID)),
	); err != nil {
		return fmt.Errorf("msg edit: %w", err)
	}

	return nil
}
