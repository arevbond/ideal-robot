-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE "scenario_actions"(
                                   "id" UUID NOT NULL,
                                   "scenario_id" UUID NOT NULL,
                                   "device_id" UUID NOT NULL,
                                   "action_type" INTEGER NOT NULL
);
ALTER TABLE
    "scenario_actions" ADD PRIMARY KEY("id");
CREATE TABLE "scenarios"(
                            "id" UUID NOT NULL,
                            "user_id" UUID NOT NULL,
                            "name" VARCHAR(255) NOT NULL,
                            "description" VARCHAR(255) NOT NULL
);
ALTER TABLE
    "scenarios" ADD PRIMARY KEY("id");
CREATE TABLE "scenario_conditions"(
                                      "id" UUID NOT NULL,
                                      "scenario_id" UUID NOT NULL,
                                      "device_id" UUID NOT NULL,
                                      "condtition_type" INTEGER NOT NULL
);
ALTER TABLE
    "scenario_conditions" ADD PRIMARY KEY("id");
CREATE TABLE "hubs"(
                       "id" UUID NOT NULL,
                       "user_id" UUID NOT NULL,
                       "name" VARCHAR(255) NOT NULL,
                       "description" VARCHAR(255) NOT NULL
);
ALTER TABLE
    "hubs" ADD PRIMARY KEY("id");
CREATE TABLE "users"(
                        "id" UUID NOT NULL,
                        "username" VARCHAR(255) NOT NULL,
                        "password_hash" VARCHAR(255) NOT NULL,
                        "email" VARCHAR(255) NOT NULL,
                        "created_at" TIMESTAMP(0) WITH TIME ZONE NOT NULL
);
ALTER TABLE
    "users" ADD PRIMARY KEY("id");
CREATE TABLE "device_commands"(
                                  "id" UUID NOT NULL,
                                  "device_id" UUID NOT NULL,
                                  "user_id" UUID NOT NULL,
                                  "created_at" TIMESTAMP(0) WITH TIME ZONE NOT NULL,
                                  "command_type" INTEGER NOT NULL
);
ALTER TABLE
    "device_commands" ADD PRIMARY KEY("id");
CREATE TABLE "device_types"(
                               "id" UUID NOT NULL,
                               "name" VARCHAR(255) NOT NULL,
                               "description" VARCHAR(255) NOT NULL,
                               "unit" BIGINT NOT NULL
);
ALTER TABLE
    "device_types" ADD PRIMARY KEY("id");
CREATE TABLE "devices"(
                          "id" UUID NOT NULL,
                          "hub_id" UUID NOT NULL,
                          "name" VARCHAR(255) NOT NULL,
                          "type" INTEGER NOT NULL,
                          "location" VARCHAR(255) NOT NULL,
                          "status" BOOLEAN NOT NULL
);
ALTER TABLE
    "devices" ADD PRIMARY KEY("id");
CREATE TABLE "device_data"(
                              "id" UUID NOT NULL,
                              "device_id" UUID NOT NULL,
                              "device_type_id" UUID NOT NULL,
                              "value" jsonb NOT NULL,
                              "recieve_at" TIMESTAMP(0) WITH TIME ZONE NOT NULL
);
ALTER TABLE
    "device_data" ADD PRIMARY KEY("id");
ALTER TABLE
    "device_commands" ADD CONSTRAINT "device_commands_device_id_foreign" FOREIGN KEY("device_id") REFERENCES "devices"("id");
ALTER TABLE
    "scenario_actions" ADD CONSTRAINT "scenario_actions_device_id_foreign" FOREIGN KEY("device_id") REFERENCES "devices"("id");
ALTER TABLE
    "scenario_conditions" ADD CONSTRAINT "scenario_conditions_scenario_id_foreign" FOREIGN KEY("scenario_id") REFERENCES "scenarios"("id");
ALTER TABLE
    "hubs" ADD CONSTRAINT "hubs_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "users"("id");
ALTER TABLE
    "device_commands" ADD CONSTRAINT "device_commands_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "users"("id");
ALTER TABLE
    "scenarios" ADD CONSTRAINT "scenarios_id_foreign" FOREIGN KEY("id") REFERENCES "scenario_actions"("id");
ALTER TABLE
    "scenarios" ADD CONSTRAINT "scenarios_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "users"("id");
ALTER TABLE
    "device_data" ADD CONSTRAINT "device_data_device_type_id_foreign" FOREIGN KEY("device_type_id") REFERENCES "device_types"("id");
ALTER TABLE
    "scenario_conditions" ADD CONSTRAINT "scenario_conditions_device_id_foreign" FOREIGN KEY("device_id") REFERENCES "devices"("id");
ALTER TABLE
    "devices" ADD CONSTRAINT "devices_hub_id_foreign" FOREIGN KEY("hub_id") REFERENCES "hubs"("id");
ALTER TABLE
    "device_data" ADD CONSTRAINT "device_data_device_id_foreign" FOREIGN KEY("device_id") REFERENCES "devices"("id");
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
