package processor

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

type ProcessorSuit struct {
	*suite.Suite
	ctrl *gomock.Controller

	mangler           *port.MockMangler
	msgSender         *port.MockMsgSender
	wordsProvider     *port.MockWordsProvider
	wordsUpdater      *port.MockWordsUpdater
	permissionChecker *port.MockPermissionChecker
	commandStorage    *port.MockCommandStorage
	banProcessor      *port.MockBanProcessor
	rankings          *port.MockRankings

	processor *processor
}

func TestProcessorSuit(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ProcessorSuit{
		Suite: new(suite.Suite),
	})
}

func (s *ProcessorSuit) SetupSuite() {
	s.ctrl = gomock.NewController(s.T())

	s.mangler = port.NewMockMangler(s.ctrl)
	s.msgSender = port.NewMockMsgSender(s.ctrl)
	s.wordsProvider = port.NewMockWordsProvider(s.ctrl)
	s.wordsUpdater = port.NewMockWordsUpdater(s.ctrl)
	s.permissionChecker = port.NewMockPermissionChecker(s.ctrl)
	s.commandStorage = port.NewMockCommandStorage(s.ctrl)
	s.banProcessor = port.NewMockBanProcessor(s.ctrl)
	s.rankings = port.NewMockRankings(s.ctrl)

	s.processor = New(
		s.mangler,
		s.msgSender,
		s.wordsProvider,
		s.wordsUpdater,
		s.permissionChecker,
		s.commandStorage,
		s.banProcessor,
		s.rankings,
	)
}

func (s *ProcessorSuit) TearDownSuite() {
}

func (s *ProcessorSuit) TearDownTest() {
}

func (s *ProcessorSuit) TearDownSubTest() {
}
