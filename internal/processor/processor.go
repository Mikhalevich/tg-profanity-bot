package processor

import (
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/cmd"
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

	cmdRouter     cmd.Router
	buttonsRouter cmd.Router
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
	}

	p.initCommandRoutes()
	p.initButtonsRoutes()

	return p
}

func (p *processor) initCommandRoutes() {
	p.cmdRouter = cmd.Router{
		cmd.GetAll: {
			Handler: p.GetAllWords,
			Perm:    cmd.Admin,
		},
	}

	if p.wordsUpdater != nil {
		p.cmdRouter[cmd.Add] = cmd.Route{
			Handler: p.AddWordCommand,
			Perm:    cmd.Admin,
		}

		p.cmdRouter[cmd.Remove] = cmd.Route{
			Handler: p.RemoveWordCommand,
			Perm:    cmd.Admin,
		}
	}
}

func (p *processor) initButtonsRoutes() {
	if p.wordsUpdater == nil {
		return
	}

	p.buttonsRouter = cmd.Router{
		cmd.Add: {
			Handler: p.AddWordCallbackQuery,
			Perm:    cmd.Admin,
		},
		cmd.Remove: {
			Handler: p.RemoveWordCallbackQuery,
			Perm:    cmd.Admin,
		},
		cmd.ViewOrginMsg: {
			Handler: p.ViewOriginMsgCallbackQuery,
			Perm:    cmd.Admin,
		},
		cmd.ViewBannedMsg: {
			Handler: p.ViewBannedMsgCallbackQuery,
			Perm:    cmd.Admin,
		},
		cmd.Unban: {
			Handler: p.UnbanCallbackQuery,
			Perm:    cmd.Admin,
		},
	}
}
