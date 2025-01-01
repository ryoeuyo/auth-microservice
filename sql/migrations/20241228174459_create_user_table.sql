-- +goose Up
CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY,
  login VARCHAR(256) NOT NULL UNIQUE,
  passHash BYTEA NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_login ON users (login);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP INDEX idx_login;
DROP TABLE users;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
