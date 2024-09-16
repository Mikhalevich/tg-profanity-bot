package router

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

type Router[T comparable] struct {
	routes map[T]Route
}

func NewRouter[T comparable]() *Router[T] {
	return &Router[T]{
		routes: make(map[T]Route),
	}
}

func (r *Router[T]) AddRoute(pattern T, handler Route) {
	r.routes[pattern] = handler
}

func (r *Router[T]) Route(pattern T) (Route, bool) {
	route, ok := r.routes[pattern]
	return route, ok
}
