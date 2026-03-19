-- +goose Up
-- +goose StatementBegin
ALTER TABLE couriers 
ADD COLUMN transport_type TEXT NOT NULL DEFAULT 'on_foot' 
CHECK (transport_type IN ('on_foot', 'scooter', 'car'));

UPDATE couriers SET transport_type = 'on_foot' WHERE transport_type IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE couriers DROP COLUMN transport_type;
-- +goose StatementEnd