//nolint:dupl
package processor

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/cbquery"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (p *processor) RemoveWordCommand(
	ctx context.Context,
	info port.MessageInfo,
	cmdArgs string,
) error {
	word := strings.TrimSpace(cmdArgs)

	return p.removeWord(
		ctx,
		info,
		word,
		port.WithButton(p.revertButton(ctx, cbquery.Add, word)),
	)
}

func (p *processor) RemoveWordCallbackQuery(
	ctx context.Context,
	info port.MessageInfo,
	word string,
) error {
	return p.removeWord(ctx, info, word, nil)
}

func (p *processor) removeWord(
	ctx context.Context,
	info port.MessageInfo,
	word string,
	options ...port.Option,
) error {
	if err := p.wordsUpdater.RemoveWord(ctx, info.ChatID.String(), word); err != nil {
		if !p.wordsUpdater.IsNothingUpdatedError(err) {
			return fmt.Errorf("remove word: %w", err)
		}

		if err := p.msgSender.Reply(ctx, info, "no such word"); err != nil {
			return fmt.Errorf("reply no such word: %w", err)
		}

		return nil
	}

	if err := p.msgSender.Reply(ctx, info, "words updated successfully", options...); err != nil {
		return fmt.Errorf("success reply: %w", err)
	}

	return nil
}
