-- +goose Up
-- +goose StatementBegin
ALTER TABLE tackle_types 
RENAME COLUMN category TO description;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tackle_types 
RENAME COLUMN description TO category;
-- +goose StatementEnd
