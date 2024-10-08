//nolint:dupl
package processor

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/cbquery"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (p *processor) AddWordCommand(ctx context.Context, info port.MessageInfo, cmdArgs string) error {
	word := strings.TrimSpace(cmdArgs)

	return p.addWord(
		ctx,
		info,
		word,
		port.WithButton(p.revertButton(ctx, cbquery.Remove, word)),
	)
}

func (p *processor) AddWordCallbackQuery(ctx context.Context, info port.MessageInfo, word string) error {
	return p.addWord(ctx, info, word, nil)
}

func (p *processor) addWord(
	ctx context.Context,
	info port.MessageInfo,
	word string,
	options ...port.Option,
) error {
	if err := p.wordsUpdater.AddWord(ctx, info.ChatID.String(), word); err != nil {
		if !p.wordsUpdater.IsNothingUpdatedError(err) {
			return fmt.Errorf("add word: %w", err)
		}

		if err := p.msgSender.Reply(ctx, info, "this word already exists"); err != nil {
			return fmt.Errorf("reply already exists: %w", err)
		}

		return nil
	}

	if err := p.msgSender.Reply(ctx, info, "words updated successfully", options...); err != nil {
		return fmt.Errorf("success reply: %w", err)
	}

	return nil
}
