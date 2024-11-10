package processor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (p *processor) TextMessage(ctx context.Context, info port.MessageInfo) error {
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
		if err := p.editBanMessage(ctx, info.UserID.String(), info); err != nil {
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
