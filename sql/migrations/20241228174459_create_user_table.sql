-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  login VARCHAR(256) NOT NULL UNIQUE,
  passHash BYTEA NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_login ON users (login);
-- +goose StatementBegin11
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP INDEX idx_login;
DROP TABLE users;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
