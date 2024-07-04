package postgres

import "context"

func (s *PostgresSuit) TestChatNotExists() {
	words, err := s.p.ChatWords(context.Background(), "chat_id")
	s.Require().EqualError(err, errChatNotExists.Error())
	s.Require().True(s.p.IsChatNotExistsError(err))
	s.Require().Nil(words)
}
