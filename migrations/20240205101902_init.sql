-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "users"(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP(0) WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS "hubs"(
   id SERIAL PRIMARY KEY NOT NULL UNIQUE,
   user_id UUID,
   FOREIGN KEY(user_id) REFERENCES "users"(id),
   name VARCHAR(255) NOT NULL,
   description VARCHAR(255)
);
CREATE TABLE IF NOT EXISTS "devices"(
    id SERIAL PRIMARY KEY NOT NULL UNIQUE,
    hub_id SERIAL,
    FOREIGN KEY (hub_id) REFERENCES "hubs"(id),
    name VARCHAR(255) NOT NULL,
    type int NOT NULL,
    location VARCHAR(255),
    status boolean default false
);
CREATE TABLE IF NOT EXISTS "devices_data"(
    id SERIAL PRIMARY KEY NOT NULL UNIQUE,
    device_id SERIAL,
    FOREIGN KEY(device_id) REFERENCES "devices"(id),
    value VARCHAR(255),
    unit VARCHAR(255),
    received_at TIMESTAMP(0) WITH TIME ZONE
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
