package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/cbquery"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (h *handler) AddWordCommand(ctx context.Context, info port.MessageInfo, cmdArgs string) error {
	word := strings.TrimSpace(cmdArgs)

	return h.addWord(
		ctx,
		info,
		word,
		func() []port.Option {
			return []port.Option{
				port.WithButton(h.revertButton(ctx, cbquery.Remove, word)),
			}
		},
	)
}

func (h *handler) AddWordCallbackQuery(ctx context.Context, info port.MessageInfo, word string) error {
	return h.addWord(
		ctx,
		info,
		word,
		nopeDelayedOption,
	)
}

type delayedOption func() []port.Option

func nopeDelayedOption() []port.Option {
	return nil
}

func (h *handler) addWord(
	ctx context.Context,
	info port.MessageInfo,
	word string,
	options delayedOption,
) error {
	if err := h.wordsUpdater.AddWord(ctx, info.ChatID.String(), word); err != nil {
		if !h.wordsUpdater.IsNothingUpdatedError(err) {
			return fmt.Errorf("add word: %w", err)
		}

		if err := h.msgSender.Reply(ctx, info, "this word already exists"); err != nil {
			return fmt.Errorf("reply already exists: %w", err)
		}

		return nil
	}

	if err := h.msgSender.Reply(ctx, info, "words updated successfully", options()...); err != nil {
		return fmt.Errorf("success reply: %w", err)
	}

	return nil
}
