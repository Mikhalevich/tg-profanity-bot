package processor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (p *processor) RequestCallbackQueryCommand(
	ctx context.Context,
	info port.MessageInfo,
	id string,
) (port.Command, bool, error) {
	command, err := p.commandStorage.Get(ctx, id)
	if err != nil {
		if !p.commandStorage.IsNotFoundError(err) {
			return port.Command{}, false, fmt.Errorf("get command from store: %w", err)
		}

		if err := p.msgSender.Reply(ctx, info, "command expired"); err != nil {
			return port.Command{}, false, fmt.Errorf("send command expired: %w", err)
		}
	}

	return command, false, nil
}
