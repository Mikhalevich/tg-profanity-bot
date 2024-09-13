package postgres

import (
	"context"
)

func (s *PostgresSuit) TestChatNotExists() {
	words, err := s.p.ChatWords(context.Background(), "chat_id")
	s.Require().EqualError(err, errChatNotExists.Error())
	s.Require().True(s.p.IsChatNotExistsError(err))
	s.Require().Nil(words)
}

func (s *PostgresSuit) TestCreateAndGetWords() {
	var (
		ctx    = context.Background()
		chatID = "chat_id"

		tests = []struct {
			Name  string
			Words []string
		}{
			{
				Name: "existing words",
				Words: []string{
					"one",
					"two",
					"three",
				},
			},
			{
				Name:  "empty slice words",
				Words: []string{},
			},
			{
				Name:  "nil slice words",
				Words: nil,
			},
		}
	)

	for _, test := range tests {
		s.Run(test.Name, func() {
			err := s.p.CreateChatWords(ctx, chatID, test.Words)
			s.Require().NoError(err)

			actualWords, err := s.p.ChatWords(ctx, chatID)
			s.Require().NoError(err)
			s.Require().ElementsMatch(test.Words, actualWords)
		})
	}
}

func (s *PostgresSuit) TestCreateAlreadyExistedChat() {
	var (
		ctx    = context.Background()
		chatID = "chat_id"
		words  = []string{
			"one",
			"two",
			"three",
		}
	)

	err := s.p.CreateChatWords(ctx, chatID, words)
	s.Require().NoError(err)

	err = s.p.CreateChatWords(ctx, chatID, words)
	s.Require().EqualError(err,
		`insert chat words: ERROR: duplicate key value violates unique constraint "chat_words_pkey" (SQLSTATE 23505)`,
	)
}

func (s *PostgresSuit) TestAddWord() {
	var (
		ctx    = context.Background()
		chatID = "chat_id"

		tests = []struct {
			Name          string
			Words         []string
			WordToAdd     string
			ExpectedWords []string
		}{
			{
				Name: "existing words",
				Words: []string{
					"one",
					"two",
					"three",
				},
				WordToAdd: "new_word",
				ExpectedWords: []string{
					"one",
					"two",
					"three",
					"new_word",
				},
			},
			{
				Name:      "empty slice words",
				Words:     []string{},
				WordToAdd: "new_word",
				ExpectedWords: []string{
					"new_word",
				},
			},
			{
				Name:      "nil slice words",
				Words:     nil,
				WordToAdd: "new_word",
				ExpectedWords: []string{
					"new_word",
				},
			},
		}
	)

	for _, test := range tests {
		s.Run(test.Name, func() {
			err := s.p.CreateChatWords(ctx, chatID, test.Words)
			s.Require().NoError(err)

			err = s.p.AddWord(ctx, chatID, test.WordToAdd)
			s.Require().NoError(err)

			actualWords, err := s.p.ChatWords(ctx, chatID)
			s.Require().NoError(err)
			s.Require().ElementsMatch(test.ExpectedWords, actualWords)
		})
	}
}

func (s *PostgresSuit) TestAddAlreadyExistingWord() {
	var (
		ctx    = context.Background()
		chatID = "chat_id"
		words  = []string{
			"one",
			"two",
			"three",
		}
	)

	err := s.p.CreateChatWords(ctx, chatID, words)
	s.Require().NoError(err)

	err = s.p.AddWord(ctx, chatID, "two")
	s.Require().EqualError(err, errNothingUpdated.Error())

	actualWords, err := s.p.ChatWords(ctx, chatID)
	s.Require().NoError(err)
	s.Require().ElementsMatch(words, actualWords)
}

func (s *PostgresSuit) TestRemoveWord() {
	var (
		ctx    = context.Background()
		chatID = "chat_id"
		words  = []string{
			"one",
			"two",
			"three",
		}
		expectedWords = []string{
			"one",
			"three",
		}
	)

	err := s.p.CreateChatWords(ctx, chatID, words)
	s.Require().NoError(err)

	err = s.p.RemoveWord(ctx, chatID, "two")
	s.Require().NoError(err)

	actualWords, err := s.p.ChatWords(ctx, chatID)
	s.Require().NoError(err)
	s.Require().ElementsMatch(expectedWords, actualWords)
}

func (s *PostgresSuit) TestRemoveNotExistingWord() {
	var (
		ctx    = context.Background()
		chatID = "chat_id"

		tests = []struct {
			Name         string
			Words        []string
			WordToRemove string
		}{
			{
				Name: "existing words",
				Words: []string{
					"one",
					"two",
					"three",
				},
				WordToRemove: "not_existing_word",
			},
			{
				Name:         "empty slice words",
				Words:        []string{},
				WordToRemove: "not_existing_word",
			},
			{
				Name:         "nil slice words",
				Words:        nil,
				WordToRemove: "not_existing_word",
			},
		}
	)

	for _, test := range tests {
		s.Run(test.Name, func() {
			err := s.p.CreateChatWords(ctx, chatID, test.Words)
			s.Require().NoError(err)

			err = s.p.RemoveWord(ctx, chatID, test.WordToRemove)
			s.Require().EqualError(err, errNothingUpdated.Error())

			actualWords, err := s.p.ChatWords(ctx, chatID)
			s.Require().NoError(err)
			s.Require().ElementsMatch(test.Words, actualWords)
		})
	}
}

func (s *PostgresSuit) TestClearWords() {
	var (
		ctx    = context.Background()
		chatID = "chat_id"

		tests = []struct {
			Name         string
			InitialWords []string
		}{
			{
				Name: "existing words",
				InitialWords: []string{
					"one",
					"two",
					"three",
				},
			},
			{
				Name:         "empty slice words",
				InitialWords: []string{},
			},
		}
	)

	for _, test := range tests {
		s.Run(test.Name, func() {
			err := s.p.CreateChatWords(ctx, chatID, test.InitialWords)
			s.Require().NoError(err)

			err = s.p.ClearWords(ctx, chatID)
			s.Require().NoError(err)

			actualWords, err := s.p.ChatWords(ctx, chatID)
			s.Require().NoError(err)
			s.Require().Len(actualWords, 0)
		})
	}
}

func (s *PostgresSuit) TestClearWordsNothingUpdatedError() {
	var (
		ctx    = context.Background()
		chatID = "chat_id"
	)

	err := s.p.ClearWords(ctx, chatID)
	s.Require().EqualError(err, errNothingUpdated.Error())
}
