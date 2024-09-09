-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS vaults(
    id SERIAL PRIMARY KEY,
    name VARCHAR (255) NOT NULL,
    content BYTEA,
    version UUID DEFAULT gen_random_uuid() NOT NULL UNIQUE,
    owner_id INTEGER REFERENCES users (id),
    UNIQUE (name, owner_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS vaults;
-- +goose StatementEnd
