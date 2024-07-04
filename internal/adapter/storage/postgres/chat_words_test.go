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
	s.Require().EqualError(err, `insert chat words: ERROR: duplicate key value violates unique constraint "chat_words_pkey" (SQLSTATE 23505)`)
}
