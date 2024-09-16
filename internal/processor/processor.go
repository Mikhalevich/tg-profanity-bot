package processor

import (
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/cbquery"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/cmd"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/router"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

type processor struct {
	mangler           port.Mangler
	msgSender         port.MsgSender
	wordsProvider     port.WordsProvider
	wordsUpdater      port.WordsUpdater
	permissionChecker port.PermissionChecker
	commandStorage    port.CommandStorage
	banProcessor      port.BanProcessor

	cmdRouter     *router.Router[cmd.CMD]
	buttonsRouter *router.Router[cbquery.CBQUERY]
}

func New(
	mangler port.Mangler,
	msgSender port.MsgSender,
	wordsProvider port.WordsProvider,
	wordsUpdater port.WordsUpdater,
	permissionChecker port.PermissionChecker,
	commandStorage port.CommandStorage,
	banProcessor port.BanProcessor,
) *processor {
	p := &processor{
		mangler:           mangler,
		msgSender:         msgSender,
		wordsProvider:     wordsProvider,
		wordsUpdater:      wordsUpdater,
		permissionChecker: permissionChecker,
		commandStorage:    commandStorage,
		banProcessor:      banProcessor,
		cmdRouter:         router.NewRouter[cmd.CMD](),
		buttonsRouter:     router.NewRouter[cbquery.CBQUERY](),
	}

	p.initCommandRoutes()
	p.initButtonsRoutes()

	return p
}

func (p *processor) initCommandRoutes() {
	p.cmdRouter.AddRoute(cmd.GetAll, router.Route{
		Handler: p.GetAllWords,
		Perm:    router.Admin,
	})

	if p.wordsUpdater != nil {
		p.cmdRouter.AddRoute(cmd.Add, router.Route{
			Handler: p.AddWordCommand,
			Perm:    router.Admin,
		})

		p.cmdRouter.AddRoute(cmd.Remove, router.Route{
			Handler: p.RemoveWordCommand,
			Perm:    router.Admin,
		})

		p.cmdRouter.AddRoute(cmd.ClearAll, router.Route{
			Handler: p.ClearWordsCommand,
			Perm:    router.Admin,
		})

		p.cmdRouter.AddRoute(cmd.RestoreDefault, router.Route{
			Handler: p.RestoreDefaultWordsCommand,
			Perm:    router.Admin,
		})
	}
}

func (p *processor) initButtonsRoutes() {
	p.buttonsRouter.AddRoute(cbquery.ViewOrginMsg, router.Route{
		Handler: p.ViewOriginMsgCallbackQuery,
		Perm:    router.Admin,
	})

	p.buttonsRouter.AddRoute(cbquery.ViewBannedMsg, router.Route{
		Handler: p.ViewBannedMsgCallbackQuery,
		Perm:    router.Admin,
	})

	p.buttonsRouter.AddRoute(cbquery.Unban, router.Route{
		Handler: p.UnbanCallbackQuery,
		Perm:    router.Admin,
	})

	if p.wordsUpdater != nil {
		p.buttonsRouter.AddRoute(cbquery.Add, router.Route{
			Handler: p.AddWordCallbackQuery,
			Perm:    router.Admin,
		})

		p.buttonsRouter.AddRoute(cbquery.Remove, router.Route{
			Handler: p.RemoveWordCallbackQuery,
			Perm:    router.Admin,
		})
	}
}
