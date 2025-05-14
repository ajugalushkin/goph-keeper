-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS vaults
    ADD file_id VARCHAR (255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS vaults
DROP COLUMN file_id;
-- +goose StatementEnd
