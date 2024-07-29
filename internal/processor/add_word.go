//nolint:dupl
package processor

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/cmd"
)

func (p *processor) AddWordCommand(ctx context.Context, chatID string, cmdArgs string, msg *tgbotapi.Message) error {
	word := strings.TrimSpace(cmdArgs)
	return p.addWord(ctx, chatID, word, msg, p.revertButton(ctx, cmd.Remove, word))
}

func (p *processor) AddWordCallbackQuery(ctx context.Context, chatID string, word string, msg *tgbotapi.Message) error {
	return p.addWord(ctx, chatID, word, msg, nil)
}

func (p *processor) addWord(
	ctx context.Context,
	chatID string,
	word string,
	msg *tgbotapi.Message,
	buttons []tgbotapi.InlineKeyboardButton,
) error {
	if err := p.wordsUpdater.AddWord(ctx, chatID, word); err != nil {
		if !p.wordsUpdater.IsNothingUpdatedError(err) {
			return fmt.Errorf("add word: %w", err)
		}

		if err := p.msgSender.Reply(ctx, msg, "this word already exists"); err != nil {
			return fmt.Errorf("reply already exists: %w", err)
		}

		return nil
	}

	if err := p.msgSender.Reply(ctx, msg, "words updated successfully", buttons...); err != nil {
		return fmt.Errorf("success reply: %w", err)
	}

	return nil
}
