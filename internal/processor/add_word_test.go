package processor

import (
	"context"
	"errors"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (s *ProcessorSuit) TestWordsUpdaterError() {
	var (
		ctx                 = context.Background()
		messageID           = 123
		chatID        int64 = 456
		userID        int64 = 789
		word                = "word"
		unexpectedErr       = errors.New("some error")
	)

	s.wordsUpdater.EXPECT().AddWord(ctx, "456", word).Return(unexpectedErr)

	s.wordsUpdater.EXPECT().IsNothingUpdatedError(unexpectedErr).Return(false)

	err := s.processor.addWord(ctx, port.MessageInfo{
		MessageID: messageID,
		ChatID:    port.NewID(chatID),
		UserID:    port.NewID(userID),
	}, word, nil)

	s.Require().Error(err, "add word: some error")
}

func (s *ProcessorSuit) TestWordAlreadyExists() {
	var (
		ctx             = context.Background()
		messageID       = 123
		chatID    int64 = 456
		userID    int64 = 789
		word            = "word"
		msgInfo         = port.MessageInfo{
			MessageID: messageID,
			ChatID:    port.NewID(chatID),
			UserID:    port.NewID(userID),
		}
		unexpectedErr = errors.New("already exists")
	)

	s.wordsUpdater.EXPECT().AddWord(ctx, "456", word).Return(unexpectedErr)

	s.wordsUpdater.EXPECT().IsNothingUpdatedError(unexpectedErr).Return(true)

	s.msgSender.EXPECT().Reply(ctx, msgInfo, "this word already exists")

	err := s.processor.addWord(ctx, msgInfo, word, nil)

	s.Require().NoError(err)
}
