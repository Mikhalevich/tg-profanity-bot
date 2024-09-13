package processor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (p *processor) ClearWordsCommand(ctx context.Context, info port.MessageInfo, cmdArgs string) error {
	if err := p.wordsUpdater.ClearWords(ctx, info.ChatID.String()); err != nil {
		if !p.wordsUpdater.IsNothingUpdatedError(err) {
			return fmt.Errorf("clear words: %w", err)
		}

		if err := p.msgSender.Reply(ctx, info, "chat does not exists"); err != nil {
			return fmt.Errorf("reply already exists: %w", err)
		}

		return nil
	}

	if err := p.msgSender.Reply(ctx, info, "words cleared successfully"); err != nil {
		return fmt.Errorf("success reply: %w", err)
	}

	return nil
}
