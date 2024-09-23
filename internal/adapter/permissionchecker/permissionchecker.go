package permissionchecker

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Mikhalevich/tg-profanity-bot/internal/app/logger"
	"github.com/Mikhalevich/tg-profanity-bot/internal/app/tracing"
)

type memberapi struct {
	api *tgbotapi.BotAPI
}

func New(api *tgbotapi.BotAPI) *memberapi {
	return &memberapi{
		api: api,
	}
}

func (m *memberapi) IsAdmin(ctx context.Context, chatID, userID int64) bool {
	_, span := tracing.StartSpan(ctx)
	defer span.End()

	chat, err := m.api.GetChat(tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: chatID,
		},
	})

	if err != nil {
		logger.FromContext(ctx).WithError(err).Warn("failed to get chat")
		return false
	}

	if chat.IsPrivate() {
		return true
	}

	member, err := m.api.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatID,
			UserID: userID,
		},
	})

	if err != nil {
		logger.FromContext(ctx).WithError(err).Warn("failed to get chat member")
		return false
	}

	return member.IsAdministrator() || member.IsCreator()
}

func (m *memberapi) UserName(ctx context.Context, chatID, userID int64) (string, error) {
	member, err := m.api.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatID,
			UserID: userID,
		},
	})

	if err != nil {
		return "", fmt.Errorf("get chat member: %w", err)
	}

	return member.User.String(), nil
}
