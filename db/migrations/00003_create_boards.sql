-- +goose Up
CREATE TABLE boards (
  id         TEXT        PRIMARY KEY DEFAULT gen_random_uuid(),
  name       TEXT        NOT NULL,
  user_id    TEXT        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS boards;
