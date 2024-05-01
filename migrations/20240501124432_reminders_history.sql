-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS "history"
(
    id SERIAL PRIMARY KEY NOT NULL UNIQUE,
    text varchar(255) NOT NULL,
    created_at timestamp,
    type int
);

CREATE TABLE IF NOT EXISTS "reminders"
(
    id SERIAL PRIMARY KEY NOT NULL UNIQUE,
    text varchar(255) NOT NULL,
    is_done boolean,
    priority int
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
