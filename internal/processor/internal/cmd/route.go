package cmd

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Permission int

const (
	Admin Permission = iota + 1
	Member
)

type Route struct {
	Handler func(ctx context.Context, chatID string, cmdArgs string, msg *tgbotapi.Message) error
	Perm    Permission
}

func (r Route) IsAdmin() bool {
	return r.Perm == Admin
}

type Router map[CMD]Route
