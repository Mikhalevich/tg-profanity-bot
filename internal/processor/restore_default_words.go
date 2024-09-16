package processor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (p *processor) RestoreDefaultWordsCommand(ctx context.Context, info port.MessageInfo, cmdArgs string) error {
	if err := p.wordsUpdater.CreateChatWords(ctx, info.ChatID.String(), p.wordsProvider.InitialWords()); err != nil {
		return fmt.Errorf("create chat words: %w", err)
	}

	if err := p.msgSender.Reply(ctx, info, "words updated successfully"); err != nil {
		return fmt.Errorf("success reply: %w", err)
	}

	return nil
}
