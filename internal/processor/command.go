package processor

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Mikhalevich/tg-profanity-bot/internal/app/tracing"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/button"
)

func (p *processor) tryProcessCommand(ctx context.Context, chatID string, msg *tgbotapi.Message) (bool, error) {
	cmd, args := extractCommand(msg.Text)
	if cmd == "" {
		return false, nil
	}

	r, ok := p.cmdRouter[cmd]
	if !ok {
		return false, nil
	}

	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	if r.IsAdmin() {
		if !p.memberChecker.IsAdmin(msg.Chat.ID, msg.From.ID) {
			return false, nil
		}
	}

	if err := r.Handler(ctx, chatID, args, msg); err != nil {
		return false, fmt.Errorf("handle command %s: %w", cmd, err)
	}

	return true, nil
}

func extractCommand(msg string) (string, string) {
	if !strings.HasPrefix(msg, "/") {
		return "", ""
	}

	cmd, args, _ := strings.Cut(msg[1:], " ")

	return cmd, args
}

func (p *processor) GetAllWords(ctx context.Context, chatID string, cmdArgs string, msg *tgbotapi.Message) error {
	words, err := p.wordsProvider.ChatWords(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get chat words: %w", err)
	}

	if err := p.msgSender.Reply(ctx, msg, strings.Join(words, "\n")); err != nil {
		return fmt.Errorf("msg reply: %w", err)
	}

	return nil
}

func (p *processor) AddWord(ctx context.Context, chatID string, cmdArgs string, msg *tgbotapi.Message) error {
	word := strings.TrimSpace(cmdArgs)
	if err := p.wordsUpdater.AddWord(ctx, chatID, word); err != nil {
		if !p.wordsUpdater.IsNothingUpdatedError(err) {
			return fmt.Errorf("add word: %w", err)
		}

		if err := p.msgSender.Reply(ctx, msg, "this word already exists"); err != nil {
			return fmt.Errorf("reply already exists: %w", err)
		}

		return nil
	}

	if err := p.msgSender.Reply(ctx, msg, "words updated successfully", makeRevertButton(word)); err != nil {
		return fmt.Errorf("success reply: %w", err)
	}

	return nil
}

func makeRevertButton(word string) []tgbotapi.InlineKeyboardButton {
	buttonInfo, err := button.ButtonCMDInfo{
		CMD:  "remove",
		Word: word,
	}.ToBase64()

	if err != nil {
		return nil
	}

	return tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("revert", buttonInfo),
	)
}

func (p *processor) RemoveWord(ctx context.Context, chatID string, cmdArgs string, msg *tgbotapi.Message) error {
	if err := p.wordsUpdater.RemoveWord(ctx, chatID, strings.TrimSpace(cmdArgs)); err != nil {
		if !p.wordsUpdater.IsNothingUpdatedError(err) {
			return fmt.Errorf("remove word: %w", err)
		}

		if err := p.msgSender.Reply(ctx, msg, "no such word"); err != nil {
			return fmt.Errorf("reply no such word: %w", err)
		}

		return nil
	}

	if err := p.msgSender.Reply(ctx, msg, "words updated successfully"); err != nil {
		return fmt.Errorf("success reply: %w", err)
	}

	return nil
}
