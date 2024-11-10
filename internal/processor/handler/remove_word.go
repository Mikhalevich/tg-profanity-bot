package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/cbquery"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (h *handler) RemoveWordCommand(
	ctx context.Context,
	info port.MessageInfo,
	cmdArgs string,
) error {
	word := strings.TrimSpace(cmdArgs)

	return h.removeWord(
		ctx,
		info,
		word,
		func() []port.Option {
			return []port.Option{
				port.WithButton(h.revertButton(ctx, cbquery.Add, word)),
			}
		},
	)
}

func (h *handler) RemoveWordCallbackQuery(
	ctx context.Context,
	info port.MessageInfo,
	word string,
) error {
	return h.removeWord(
		ctx,
		info,
		word,
		nopeDelayedOption,
	)
}

func (h *handler) removeWord(
	ctx context.Context,
	info port.MessageInfo,
	word string,
	options delayedOption,
) error {
	if err := h.wordsUpdater.RemoveWord(ctx, info.ChatID.String(), word); err != nil {
		if !h.wordsUpdater.IsNothingUpdatedError(err) {
			return fmt.Errorf("remove word: %w", err)
		}

		if err := h.msgSender.Reply(ctx, info, "no such word"); err != nil {
			return fmt.Errorf("reply no such word: %w", err)
		}

		return nil
	}

	if err := h.msgSender.Reply(ctx, info, "words updated successfully", options()...); err != nil {
		return fmt.Errorf("success reply: %w", err)
	}

	return nil
}
