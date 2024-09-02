package cmd

import (
	"context"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

type Permission int

const (
	Admin Permission = iota + 1
	Member
)

type Route struct {
	Handler func(ctx context.Context, info port.MessageInfo, cmdArgs string) error
	Perm    Permission
}

func (r Route) IsAdmin() bool {
	return r.Perm == Admin
}

type Router map[CMD]Route
