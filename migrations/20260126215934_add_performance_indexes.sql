-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_delivery_order_id ON delivery(order_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_delivery_order_id;
-- +goose StatementEnd