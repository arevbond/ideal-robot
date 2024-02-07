-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "users"(
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "username" VARCHAR(255) NOT NULL UNIQUE,
    "password_hash" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP(0) WITH TIME ZONE NOT NULL
);

CREATE TABLE "hubs"(
   id SERIAL PRIMARY KEY NOT NULL UNIQUE,
   "user_id" UUID NOT NULL,
   FOREIGN KEY("user_id") REFERENCES "users"("id"),
   "name" VARCHAR(255) NOT NULL,
   "description" VARCHAR(255) NOT NULL
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
