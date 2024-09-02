package port

import (
	"context"
)

type Button struct {
	Caption string
	Data    string
}

type MsgSender interface {
	Reply(ctx context.Context, originMsgInfo MessageInfo, msg string, buttons ...*Button) error
	Edit(ctx context.Context, originMsgInfo MessageInfo, msg string, buttons ...*Button) error
}
