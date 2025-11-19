-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  id serial PRIMARY KEY,
  login varchar(16) NOT NULL UNIQUE,
  password_hash text NOT NULL
);

CREATE TABLE tasks (
  id serial PRIMARY KEY,
  user_id integer NOT NULL,
  task text NOT NULL,
  completed boolean NOT NULL DEFAULT false,
  created_at timestamp NOT NULL DEFAULT NOW(),
  completed_at timestamp,
  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
