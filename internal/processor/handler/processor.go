package handler

import (
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

type handler struct {
	mangler           port.Mangler
	msgSender         port.MsgSender
	wordsProvider     port.WordsProvider
	wordsUpdater      port.WordsUpdater
	permissionChecker port.PermissionChecker
	commandStorage    port.CommandStorage
	banProcessor      port.BanProcessor
	rankings          port.Rankings
}

func New(
	mangler port.Mangler,
	msgSender port.MsgSender,
	wordsProvider port.WordsProvider,
	wordsUpdater port.WordsUpdater,
	permissionChecker port.PermissionChecker,
	commandStorage port.CommandStorage,
	banProcessor port.BanProcessor,
	rankings port.Rankings,
) *handler {
	h := &handler{
		mangler:           mangler,
		msgSender:         msgSender,
		wordsProvider:     wordsProvider,
		wordsUpdater:      wordsUpdater,
		permissionChecker: permissionChecker,
		commandStorage:    commandStorage,
		banProcessor:      banProcessor,
		rankings:          rankings,
	}

	return h
}
