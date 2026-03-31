-- +goose Up
CREATE TABLE boards (
  id         TEXT        PRIMARY KEY DEFAULT gen_random_uuid(),
  name       TEXT        NOT NULL,
  user_id    TEXT        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_boards_user_id ON boards(user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_boards_user_id;
DROP TABLE IF EXISTS boards;
