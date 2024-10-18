package processor

import (
	"context"
	"errors"

	"github.com/golang/mock/gomock"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/cbquery"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (s *ProcessorSuit) TestCommandWordsUpdaterError() {
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

	err := s.processor.AddWordCommand(ctx, port.MessageInfo{
		MessageID: messageID,
		ChatID:    port.NewID(chatID),
		UserID:    port.NewID(userID),
	}, word)

	s.Require().EqualError(err, "add word: some error")
}

func (s *ProcessorSuit) TestCallbackQueryWordsUpdaterError() {
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

	err := s.processor.AddWordCallbackQuery(ctx, port.MessageInfo{
		MessageID: messageID,
		ChatID:    port.NewID(chatID),
		UserID:    port.NewID(userID),
	}, word)

	s.Require().EqualError(err, "add word: some error")
}

func (s *ProcessorSuit) TestCommandWordAlreadyExists() {
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

	err := s.processor.AddWordCommand(ctx, msgInfo, word)

	s.Require().NoError(err)
}

func (s *ProcessorSuit) TestCommandWordAlreadyExistsReplyError() {
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

	s.msgSender.EXPECT().Reply(ctx, msgInfo, "this word already exists").Return(errors.New("reply error"))

	err := s.processor.AddWordCommand(ctx, msgInfo, word)

	s.Require().EqualError(err, "reply already exists: reply error")
}

func (s *ProcessorSuit) TestCommandWordsUpdatedSuccessfully() {
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
	)

	s.commandStorage.EXPECT().Set(ctx, gomock.Any(), port.Command{
		CMD:     cbquery.Remove.String(),
		Payload: []byte(word),
	})

	s.wordsUpdater.EXPECT().AddWord(ctx, "456", word).Return(nil)

	s.msgSender.EXPECT().Reply(ctx, msgInfo, "words updated successfully", gomock.Any())

	err := s.processor.AddWordCommand(ctx, msgInfo, word)

	s.Require().NoError(err)
}

func (s *ProcessorSuit) TestCallbackQueryWordsUpdatedSuccessfully() {
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
	)

	s.wordsUpdater.EXPECT().AddWord(ctx, "456", word).Return(nil)

	s.msgSender.EXPECT().Reply(ctx, msgInfo, "words updated successfully", gomock.Any())

	err := s.processor.AddWordCallbackQuery(ctx, msgInfo, word)

	s.Require().NoError(err)
}
