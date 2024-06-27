package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	errChatNotExists = errors.New("chat is not exists")
)

func (p *Postgres) ChatWords(ctx context.Context, chatID string) ([]string, error) {
	query, args, err := sqlx.Named(`
		SELECT words
		FROM chat_words
		WHERE chat_id = :chat_id
	`,
		map[string]any{
			"chat_id": chatID,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("sqlx named: %w", err)
	}

	var jsonb string
	if err := sqlx.GetContext(ctx, p.db, &jsonb, p.db.Rebind(query), args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errChatNotExists
		}

		return nil, fmt.Errorf("select jsonb payload: %w", err)
	}

	var words []string
	if err := json.Unmarshal([]byte(jsonb), &words); err != nil {
		return nil, fmt.Errorf("unmarshal words: %w", err)
	}

	return words, nil
}

func (p *Postgres) IsChatNotExistsError(err error) bool {
	return errors.Is(err, errChatNotExists)
}

func (p *Postgres) CreateChatWords(ctx context.Context, chatID string, words []string) error {
	payload, err := json.Marshal(words)
	if err != nil {
		return fmt.Errorf("marshal words: %w", err)
	}

	res, err := p.db.NamedExecContext(
		ctx,
		"INSERT INTO chat_words(chat_id, words) VALUES(:chat_id, :words)",
		map[string]any{
			"chat_id": chatID,
			"words":   payload,
		})

	if err != nil {
		return fmt.Errorf("insert chat words: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}

	if rows == 0 {
		//nolint:err113
		return errors.New("no words inserted")
	}

	return nil
}
