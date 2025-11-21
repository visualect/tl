-- +goose Up
-- +goose StatementBegin
ALTER TABLE tasks RENAME COLUMN completed_at TO updated_at;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tasks RENAME COLUMN updated_at TO completed_at;
-- +goose StatementEnd
