-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE chat_words(
    chat_id TEXT PRIMARY KEY,
    words JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE chat_words;