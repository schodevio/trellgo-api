-- +goose Up
CREATE TABLE lists (
  id         TEXT        PRIMARY KEY DEFAULT gen_random_uuid(),
  name       TEXT        NOT NULL,
  board_id   TEXT        NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
  position   INT         NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_lists_board_id ON lists(board_id);

-- +goose Down
DROP INDEX IF EXISTS idx_lists_board_id;
DROP TABLE IF EXISTS lists;
