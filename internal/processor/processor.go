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

		p.cmdRouter[cmd.Clear] = cmd.Route{
			Handler: p.ClearWordsCommand,
			Perm:    cmd.Admin,
		}

		p.cmdRouter[cmd.RestoreDefault] = cmd.Route{
			Handler: p.RestoreDefaultWordsCommand,
			Perm:    cmd.Admin,
		}
	}
}

func (p *processor) initButtonsRoutes() {
	p.buttonsRouter = cmd.Router{
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

	if p.wordsUpdater != nil {
		p.buttonsRouter[cmd.Add] = cmd.Route{
			Handler: p.AddWordCallbackQuery,
			Perm:    cmd.Admin,
		}

		p.buttonsRouter[cmd.Remove] = cmd.Route{
			Handler: p.RemoveWordCallbackQuery,
			Perm:    cmd.Admin,
		}
	}
}
