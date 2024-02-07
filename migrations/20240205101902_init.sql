-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- CREATE TABLE "scenario_actions"(
-- id SERIAL PRIMARY KEY NOT NULL UNIQUE,
-- "scenario_id" SERIAL NOT NULL,
-- "device_id" SERIAL NOT NULL,
-- "action_type" INTEGER NOT NULL
-- );
-- CREATE TABLE "scenarios"(
-- id SERIAL PRIMARY KEY NOT NULL UNIQUE,
-- "user_id" UUID NOT NULL,
-- "name" VARCHAR(255) NOT NULL,
-- "description" VARCHAR(255) NOT NULL
-- );
-- CREATE TABLE "scenario_conditions"(
-- id SERIAL PRIMARY KEY NOT NULL UNIQUE,
-- "scenario_id" SERIAL NOT NULL,
-- "device_id" SERIAL NOT NULL,
-- "condtition_type" INTEGER NOT NULL
-- );
CREATE TABLE "hubs"(
id SERIAL PRIMARY KEY NOT NULL UNIQUE,
"user_id" UUID NOT NULL,
"name" VARCHAR(255) NOT NULL,
"description" VARCHAR(255) NOT NULL
);
CREATE TABLE "users"(
"id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
"username" VARCHAR(255) NOT NULL UNIQUE,
"password_hash" VARCHAR(255) NOT NULL,
"email" VARCHAR(255) NOT NULL,
"created_at" TIMESTAMP(0) WITH TIME ZONE NOT NULL
);
-- CREATE TABLE "device_commands"(
-- id SERIAL PRIMARY KEY NOT NULL UNIQUE,
-- "device_id" SERIAL NOT NULL,
-- "user_id" UUID NOT NULL,
-- "created_at" TIMESTAMP(0) WITH TIME ZONE NOT NULL,
-- "command_type" INTEGER NOT NULL
-- );
-- CREATE TABLE "device_types"(
-- id SERIAL PRIMARY KEY NOT NULL UNIQUE,
-- "name" VARCHAR(255) NOT NULL,
-- "description" VARCHAR(255) NOT NULL,
-- "unit" BIGINT NOT NULL
-- );
-- CREATE TABLE "devices"(
-- id SERIAL PRIMARY KEY NOT NULL UNIQUE,
-- "hub_id" SERIAL NOT NULL,
-- "name" VARCHAR(255) NOT NULL,
-- "type" INTEGER NOT NULL,
-- "location" VARCHAR(255) NOT NULL,
-- "status" BOOLEAN NOT NULL
-- );
-- CREATE TABLE "device_data"(
-- id SERIAL PRIMARY KEY NOT NULL UNIQUE,
-- "device_id" SERIAL NOT NULL,
-- "device_type_id" INTEGER NOT NULL,
-- "value" jsonb NOT NULL,
-- "recieve_at" TIMESTAMP(0) WITH TIME ZONE NOT NULL
-- );
-- ALTER TABLE
-- "device_commands" ADD CONSTRAINT "device_commands_device_id_foreign" FOREIGN KEY("device_id") REFERENCES "devices"("id");
-- ALTER TABLE
-- "scenario_actions" ADD CONSTRAINT "scenario_actions_device_id_foreign" FOREIGN KEY("device_id") REFERENCES "devices"("id");
-- ALTER TABLE
-- "scenario_conditions" ADD CONSTRAINT "scenario_conditions_scenario_id_foreign" FOREIGN KEY("scenario_id") REFERENCES "scenarios"("id");
ALTER TABLE
"hubs" ADD CONSTRAINT "hubs_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "users"("id");
-- ALTER TABLE
-- "device_commands" ADD CONSTRAINT "device_commands_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "users"("id");
-- ALTER TABLE
-- "scenarios" ADD CONSTRAINT "scenarios_id_foreign" FOREIGN KEY("id") REFERENCES "scenario_actions"("id");
-- ALTER TABLE
-- "scenarios" ADD CONSTRAINT "scenarios_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "users"("id");
-- ALTER TABLE
-- "device_data" ADD CONSTRAINT "device_data_device_type_id_foreign" FOREIGN KEY("device_type_id") REFERENCES "device_types"("id");
-- ALTER TABLE
-- "scenario_conditions" ADD CONSTRAINT "scenario_conditions_device_id_foreign" FOREIGN KEY("device_id") REFERENCES "devices"("id");
-- ALTER TABLE
-- "devices" ADD CONSTRAINT "devices_hub_id_foreign" FOREIGN KEY("hub_id") REFERENCES "hubs"("id");
-- ALTER TABLE
-- "device_data" ADD CONSTRAINT "device_data_device_id_foreign" FOREIGN KEY("device_id") REFERENCES "devices"("id");
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
