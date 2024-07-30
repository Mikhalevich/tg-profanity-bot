//nolint:dupl
package processor

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/cmd"
)

func (p *processor) RemoveWordCommand(
	ctx context.Context,
	chatID string,
	cmdArgs string,
	msg *tgbotapi.Message,
) error {
	word := strings.TrimSpace(cmdArgs)
	return p.removeWord(ctx, chatID, word, msg, p.revertButton(ctx, cmd.Add, word))
}

func (p *processor) RemoveWordCallbackQuery(
	ctx context.Context,
	chatID string,
	word string,
	msg *tgbotapi.Message,
) error {
	return p.removeWord(ctx, chatID, word, msg, nil)
}

func (p *processor) removeWord(
	ctx context.Context,
	chatID string,
	word string,
	msg *tgbotapi.Message,
	buttons []tgbotapi.InlineKeyboardButton,
) error {
	if err := p.wordsUpdater.RemoveWord(ctx, chatID, word); err != nil {
		if !p.wordsUpdater.IsNothingUpdatedError(err) {
			return fmt.Errorf("remove word: %w", err)
		}

		if err := p.msgSender.Reply(ctx, msg, "no such word"); err != nil {
			return fmt.Errorf("reply no such word: %w", err)
		}

		return nil
	}

	if err := p.msgSender.Reply(ctx, msg, "words updated successfully", buttons...); err != nil {
		return fmt.Errorf("success reply: %w", err)
	}

	return nil
}
