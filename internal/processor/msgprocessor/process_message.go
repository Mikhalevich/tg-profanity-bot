package msgprocessor

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mikhalevich/tg-profanity-bot/internal/infra/tracing"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/cmd"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (m *msgprocessor) ProcessMessage(ctx context.Context, info port.MessageInfo) error {
	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	isBanned, err := m.handler.ProcessBannedMessage(ctx, info)
	if err != nil {
		return fmt.Errorf("is message allowed: %w", err)
	}

	if isBanned {
		return nil
	}

	isProcessed, err := m.tryProcessCommand(ctx, info)
	if err != nil {
		return fmt.Errorf("process command: %w", err)
	}

	if isProcessed {
		return nil
	}

	if err := m.handler.TextMessage(ctx, info); err != nil {
		return fmt.Errorf("process text message: %w", err)
	}

	return nil
}

func (m *msgprocessor) tryProcessCommand(ctx context.Context, info port.MessageInfo) (bool, error) {
	command, args := extractCommand(info.Text)
	if command == "" {
		return false, nil
	}

	r, ok := m.cmdRouter.Route(command)
	if !ok {
		return false, nil
	}

	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	if r.IsAdmin() {
		if !m.permissionChecker.IsAdmin(ctx, info.ChatID.Int64(), info.UserID.Int64()) {
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
