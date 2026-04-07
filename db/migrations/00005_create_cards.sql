-- +goose Up
CREATE TABLE cards (
  id          TEXT        PRIMARY KEY DEFAULT gen_random_uuid(),
  list_id     TEXT        NOT NULL REFERENCES lists(id) ON DELETE CASCADE,
  position    INT         NOT NULL,
  title       TEXT        NOT NULL,
  description TEXT,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_cards_list_id ON cards(list_id);

-- +goose Down
DROP INDEX IF EXISTS idx_cards_list_id;
DROP TABLE IF EXISTS cards;
