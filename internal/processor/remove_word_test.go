package processor

import (
	"context"
	"errors"

	"github.com/golang/mock/gomock"

	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/internal/cbquery"
	"github.com/Mikhalevich/tg-profanity-bot/internal/processor/port"
)

func (s *ProcessorSuit) TestRemoveWordCommandWordsUpdaterError() {
	var (
		ctx                 = context.Background()
		messageID           = 987
		chatID        int64 = 654
		userID        int64 = 321
		word                = "word_to_remove"
		unexpectedErr       = errors.New("some error")
	)

	s.commandStorage.EXPECT().Set(ctx, gomock.Any(), port.Command{
		CMD:     cbquery.Add.String(),
		Payload: []byte(word),
	})

	s.wordsUpdater.EXPECT().RemoveWord(ctx, "654", word).Return(unexpectedErr)

	s.wordsUpdater.EXPECT().IsNothingUpdatedError(unexpectedErr).Return(false)

	err := s.processor.RemoveWordCommand(ctx, port.MessageInfo{
		MessageID: messageID,
		ChatID:    port.NewID(chatID),
		UserID:    port.NewID(userID),
	}, word)

	s.Require().EqualError(err, "remove word: some error")
}

func (s *ProcessorSuit) TestRemoveWordCallbackQueryWordsUpdaterError() {
	var (
		ctx                 = context.Background()
		messageID           = 987
		chatID        int64 = 654
		userID        int64 = 321
		word                = "word_to_remove"
		unexpectedErr       = errors.New("some error")
	)

	s.wordsUpdater.EXPECT().RemoveWord(ctx, "654", word).Return(unexpectedErr)

	s.wordsUpdater.EXPECT().IsNothingUpdatedError(unexpectedErr).Return(false)

	err := s.processor.RemoveWordCallbackQuery(ctx, port.MessageInfo{
		MessageID: messageID,
		ChatID:    port.NewID(chatID),
		UserID:    port.NewID(userID),
	}, word)

	s.Require().EqualError(err, "remove word: some error")
}

func (s *ProcessorSuit) TestRemoveWordCommandAlreadyExistsError() {
	var (
		ctx             = context.Background()
		messageID       = 987
		chatID    int64 = 654
		userID    int64 = 321
		msgInfo         = port.MessageInfo{
			MessageID: messageID,
			ChatID:    port.NewID(chatID),
			UserID:    port.NewID(userID),
		}
		word          = "word_to_remove"
		unexpectedErr = errors.New("some error")
	)

	s.commandStorage.EXPECT().Set(ctx, gomock.Any(), port.Command{
		CMD:     cbquery.Add.String(),
		Payload: []byte(word),
	})

	s.wordsUpdater.EXPECT().RemoveWord(ctx, "654", word).Return(unexpectedErr)

	s.wordsUpdater.EXPECT().IsNothingUpdatedError(unexpectedErr).Return(true)

	s.msgSender.EXPECT().Reply(ctx, msgInfo, "no such word").Return(nil)

	err := s.processor.RemoveWordCommand(ctx, msgInfo, word)

	s.Require().NoError(err)
}

func (s *ProcessorSuit) TestRemoveWordCallbackQueryAlreadyExistsError() {
	var (
		ctx             = context.Background()
		messageID       = 987
		chatID    int64 = 654
		userID    int64 = 321
		msgInfo         = port.MessageInfo{
			MessageID: messageID,
			ChatID:    port.NewID(chatID),
			UserID:    port.NewID(userID),
		}
		word          = "word_to_remove"
		unexpectedErr = errors.New("some error")
	)

	s.wordsUpdater.EXPECT().RemoveWord(ctx, "654", word).Return(unexpectedErr)

	s.wordsUpdater.EXPECT().IsNothingUpdatedError(unexpectedErr).Return(true)

	s.msgSender.EXPECT().Reply(ctx, msgInfo, "no such word").Return(nil)

	err := s.processor.RemoveWordCallbackQuery(ctx, msgInfo, word)

	s.Require().NoError(err)
}

func (s *ProcessorSuit) TestRemoveWordCommandAlreadyExistsReplyError() {
	var (
		ctx             = context.Background()
		messageID       = 987
		chatID    int64 = 654
		userID    int64 = 321
		msgInfo         = port.MessageInfo{
			MessageID: messageID,
			ChatID:    port.NewID(chatID),
			UserID:    port.NewID(userID),
		}
		word          = "word_to_remove"
		unexpectedErr = errors.New("some error")
	)

	s.commandStorage.EXPECT().Set(ctx, gomock.Any(), port.Command{
		CMD:     cbquery.Add.String(),
		Payload: []byte(word),
	})

	s.wordsUpdater.EXPECT().RemoveWord(ctx, "654", word).Return(unexpectedErr)

	s.wordsUpdater.EXPECT().IsNothingUpdatedError(unexpectedErr).Return(true)

	s.msgSender.EXPECT().Reply(ctx, msgInfo, "no such word").Return(errors.New("some reply error"))

	err := s.processor.RemoveWordCommand(ctx, msgInfo, word)

	s.Require().EqualError(err, "reply no such word: some reply error")
}

func (s *ProcessorSuit) TestRemoveWordCallbackQueryAlreadyExistsReplyError() {
	var (
		ctx             = context.Background()
		messageID       = 987
		chatID    int64 = 654
		userID    int64 = 321
		msgInfo         = port.MessageInfo{
			MessageID: messageID,
			ChatID:    port.NewID(chatID),
			UserID:    port.NewID(userID),
		}
		word          = "word_to_remove"
		unexpectedErr = errors.New("some error")
	)

	s.wordsUpdater.EXPECT().RemoveWord(ctx, "654", word).Return(unexpectedErr)

	s.wordsUpdater.EXPECT().IsNothingUpdatedError(unexpectedErr).Return(true)

	s.msgSender.EXPECT().Reply(ctx, msgInfo, "no such word").Return(errors.New("some reply error"))

	err := s.processor.RemoveWordCallbackQuery(ctx, msgInfo, word)

	s.Require().EqualError(err, "reply no such word: some reply error")
}

func (s *ProcessorSuit) TestRemoveWordCommandReplyError() {
	var (
		ctx             = context.Background()
		messageID       = 987
		chatID    int64 = 654
		userID    int64 = 321
		msgInfo         = port.MessageInfo{
			MessageID: messageID,
			ChatID:    port.NewID(chatID),
			UserID:    port.NewID(userID),
		}
		word = "word_to_remove"
	)

	s.commandStorage.EXPECT().Set(ctx, gomock.Any(), port.Command{
		CMD:     cbquery.Add.String(),
		Payload: []byte(word),
	})

	s.wordsUpdater.EXPECT().RemoveWord(ctx, "654", word).Return(nil)

	s.msgSender.EXPECT().
		Reply(ctx, msgInfo, "words updated successfully", gomock.Any()).
		Return(errors.New("some reply error"))

	err := s.processor.RemoveWordCommand(ctx, msgInfo, word)

	s.Require().EqualError(err, "success reply: some reply error")
}

func (s *ProcessorSuit) TestRemoveWordCallbackQueryReplyError() {
	var (
		ctx             = context.Background()
		messageID       = 987
		chatID    int64 = 654
		userID    int64 = 321
		msgInfo         = port.MessageInfo{
			MessageID: messageID,
			ChatID:    port.NewID(chatID),
			UserID:    port.NewID(userID),
		}
		word = "word_to_remove"
	)

	s.wordsUpdater.EXPECT().RemoveWord(ctx, "654", word).Return(nil)

	s.msgSender.EXPECT().
		Reply(ctx, msgInfo, "words updated successfully", nil).
		Return(errors.New("some reply error"))

	err := s.processor.RemoveWordCallbackQuery(ctx, msgInfo, word)

	s.Require().EqualError(err, "success reply: some reply error")
}

func (s *ProcessorSuit) TestRemoveWordCommandSuccessReply() {
	var (
		ctx             = context.Background()
		messageID       = 987
		chatID    int64 = 654
		userID    int64 = 321
		msgInfo         = port.MessageInfo{
			MessageID: messageID,
			ChatID:    port.NewID(chatID),
			UserID:    port.NewID(userID),
		}
		word = "word_to_remove"
	)

	s.commandStorage.EXPECT().Set(ctx, gomock.Any(), port.Command{
		CMD:     cbquery.Add.String(),
		Payload: []byte(word),
	})

	s.wordsUpdater.EXPECT().RemoveWord(ctx, "654", word).Return(nil)

	s.msgSender.EXPECT().
		Reply(ctx, msgInfo, "words updated successfully", gomock.Any()).
		Return(nil)

	err := s.processor.RemoveWordCommand(ctx, msgInfo, word)

	s.Require().NoError(err)
}

func (s *ProcessorSuit) TestRemoveWordCallbackQuerySuccessReply() {
	var (
		ctx             = context.Background()
		messageID       = 987
		chatID    int64 = 654
		userID    int64 = 321
		msgInfo         = port.MessageInfo{
			MessageID: messageID,
			ChatID:    port.NewID(chatID),
			UserID:    port.NewID(userID),
		}
		word = "word_to_remove"
	)

	s.wordsUpdater.EXPECT().RemoveWord(ctx, "654", word).Return(nil)

	s.msgSender.EXPECT().
		Reply(ctx, msgInfo, "words updated successfully", gomock.Any()).
		Return(nil)

	err := s.processor.RemoveWordCallbackQuery(ctx, msgInfo, word)

	s.Require().NoError(err)
}
