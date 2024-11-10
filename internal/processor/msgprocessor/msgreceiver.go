package msgprocessor

import (
	"context"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/cbquery"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/cmd"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/router"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

//nolint:interfacebloat
type MsgHandler interface {
	AddWordCommand(ctx context.Context, info port.MessageInfo, cmdArgs string) error
	AddWordCallbackQuery(ctx context.Context, info port.MessageInfo, word string) error
	ClearWordsCommand(ctx context.Context, info port.MessageInfo, cmdArgs string) error
	RestoreDefaultWordsCommand(ctx context.Context, info port.MessageInfo, cmdArgs string) error
	GetAllWords(ctx context.Context, info port.MessageInfo, cmdArgs string) error
	Rankings(ctx context.Context, info port.MessageInfo, monthArg string) error
	RemoveWordCommand(ctx context.Context, info port.MessageInfo, cmdArgs string) error
	RemoveWordCallbackQuery(ctx context.Context, info port.MessageInfo, word string) error
	UnbanCallbackQuery(ctx context.Context, info port.MessageInfo, userID string) error
	ViewBannedMsgCallbackQuery(ctx context.Context, info port.MessageInfo, msgText string) error
	ViewOriginMsgCallbackQuery(ctx context.Context, info port.MessageInfo, originMsgText string) error
	TextMessage(ctx context.Context, info port.MessageInfo) error

	ProcessBannedMessage(ctx context.Context, info port.MessageInfo) (bool, error)
	RequestCallbackQueryCommand(ctx context.Context, info port.MessageInfo, id string) (port.Command, bool, error)
	ReplyTextMessage(ctx context.Context, originInfo port.MessageInfo, text string) error
}

type msgprocessor struct {
	handler           MsgHandler
	permissionChecker port.PermissionChecker
	cmdRouter         *router.Router[cmd.CMD]
	buttonsRouter     *router.Router[cbquery.CBQUERY]
}

func NewMsgProcessor(
	handler MsgHandler,
	permissionChecker port.PermissionChecker,
	isWordsModificationsSupported bool,
	isRankingsSupported bool,
) *msgprocessor {
	m := &msgprocessor{
		handler:           handler,
		permissionChecker: permissionChecker,
		cmdRouter:         router.NewRouter[cmd.CMD](),
		buttonsRouter:     router.NewRouter[cbquery.CBQUERY](),
	}

	m.initCommandRoutes(isWordsModificationsSupported, isRankingsSupported)
	m.initButtonsRoutes(isWordsModificationsSupported)

	return m
}

func (m *msgprocessor) initCommandRoutes(
	isWordsModificationsSupported bool,
	isRankingsSupported bool,
) {
	m.cmdRouter.AddRoute(cmd.GetAll, router.Route{
		Handler: m.handler.GetAllWords,
		Perm:    router.Admin,
	})

	if isWordsModificationsSupported {
		m.cmdRouter.AddRoute(cmd.Add, router.Route{
			Handler: m.handler.AddWordCommand,
			Perm:    router.Admin,
		})

		m.cmdRouter.AddRoute(cmd.Remove, router.Route{
			Handler: m.handler.RemoveWordCommand,
			Perm:    router.Admin,
		})

		m.cmdRouter.AddRoute(cmd.ClearAll, router.Route{
			Handler: m.handler.ClearWordsCommand,
			Perm:    router.Admin,
		})

		m.cmdRouter.AddRoute(cmd.RestoreDefault, router.Route{
			Handler: m.handler.RestoreDefaultWordsCommand,
			Perm:    router.Admin,
		})
	}

	if isRankingsSupported {
		m.cmdRouter.AddRoute(cmd.Rankings, router.Route{
			Handler: m.handler.Rankings,
			Perm:    router.Member,
		})
	}
}

func (m *msgprocessor) initButtonsRoutes(isWordsModificationsSupported bool) {
	m.buttonsRouter.AddRoute(cbquery.ViewOrginMsg, router.Route{
		Handler: m.handler.ViewOriginMsgCallbackQuery,
		Perm:    router.Admin,
	})

	m.buttonsRouter.AddRoute(cbquery.ViewBannedMsg, router.Route{
		Handler: m.handler.ViewBannedMsgCallbackQuery,
		Perm:    router.Admin,
	})

	m.buttonsRouter.AddRoute(cbquery.Unban, router.Route{
		Handler: m.handler.UnbanCallbackQuery,
		Perm:    router.Admin,
	})

	if isWordsModificationsSupported {
		m.buttonsRouter.AddRoute(cbquery.Add, router.Route{
			Handler: m.handler.AddWordCallbackQuery,
			Perm:    router.Admin,
		})

		m.buttonsRouter.AddRoute(cbquery.Remove, router.Route{
			Handler: m.handler.RemoveWordCallbackQuery,
			Perm:    router.Admin,
		})
	}
}
