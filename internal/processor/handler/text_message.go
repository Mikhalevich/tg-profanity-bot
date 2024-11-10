package handler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (h *handler) TextMessage(ctx context.Context, info port.MessageInfo) error {
	mangledMsg, err := h.mangler.Mangle(ctx, info.ChatID.String(), info.Text)
	if err != nil {
		return fmt.Errorf("replace msg: %w", err)
	}

	if mangledMsg == info.Text {
		return nil
	}

	if err := h.updateRankingScore(ctx, info); err != nil {
		return fmt.Errorf("update ranking score: %w", err)
	}

	isBanned, err := h.banProcessor.AddViolation(ctx, makeBanID(info.ChatID.String(), info.UserID.String()))
	if err != nil {
		return fmt.Errorf("add violation: %w", err)
	}

	if isBanned {
		if err := h.editBanMessage(ctx, info.UserID.String(), info); err != nil {
			return fmt.Errorf("process ban: %w", err)
		}

		return nil
	}

	if err := h.msgSender.Edit(
		ctx,
		info,
		mangledMsg,
		port.WithButton(h.viewOriginMsgButton(ctx, info.Text)),
	); err != nil {
		return fmt.Errorf("msg edit: %w", err)
	}

	return nil
}

func (h *handler) updateRankingScore(ctx context.Context, info port.MessageInfo) error {
	if err := h.rankings.AddScore(
		ctx,
		makeCurrentMonthRankingKey(info.ChatID.String()),
		info.UserID.String(),
	); err != nil {
		return fmt.Errorf("add rankings score: %w", err)
	}

	return nil
}
