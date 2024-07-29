package processor

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Mikhalevich/tg-profanity-bot/internal/app/tracing"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/button"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/cmd"
)

func (p *processor) tryProcessCommand(ctx context.Context, chatID string, msg *tgbotapi.Message) (bool, error) {
	command, args := extractCommand(msg.Text)
	if command == "" {
		return false, nil
	}

	r, ok := p.cmdRouter[command]
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
		return false, fmt.Errorf("handle command %s: %w", command.String(), err)
	}

	return true, nil
}

func extractCommand(msg string) (cmd.CMD, string) {
	if !strings.HasPrefix(msg, "/") {
		return "", ""
	}

	command, args, _ := strings.Cut(msg[1:], " ")

	return cmd.CMD(command), args
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

func revertButton(cmd, word string) []tgbotapi.InlineKeyboardButton {
	buttonInfo, err := button.ButtonCMDInfo{
		CMD:     cmd,
		Payload: []byte(word),
	}.ToBase64()

	if err != nil {
		return nil
	}

	return tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("revert", buttonInfo),
	)
}
