package processor

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mikhalevich/tg-profanity-bot/internal/app/tracing"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/cmd"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (p *processor) tryProcessCommand(ctx context.Context, info port.MessageInfo) (bool, error) {
	command, args := extractCommand(info.Text)
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
		if !p.permissionChecker.IsAdmin(ctx, info.ChatID.Int64(), info.UserID.Int64()) {
			return false, nil
		}
	}

	if err := r.Handler(ctx, info, args); err != nil {
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

func (p *processor) GetAllWords(ctx context.Context, info port.MessageInfo, cmdArgs string) error {
	words, err := p.wordsProvider.ChatWords(ctx, info.ChatID.String())
	if err != nil {
		return fmt.Errorf("get chat words: %w", err)
	}

	if err := p.msgSender.Reply(ctx, info, strings.Join(words, "\n")); err != nil {
		return fmt.Errorf("msg reply: %w", err)
	}

	return nil
}
