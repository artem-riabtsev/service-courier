-- +goose Up
-- +goose StatementBegin

CREATE INDEX IF NOT EXISTS idx_couriers_status ON couriers(status);
CREATE INDEX IF NOT EXISTS idx_couriers_phone ON couriers(phone);
CREATE INDEX IF NOT EXISTS idx_couriers_status_available ON couriers(status) WHERE status = 'available';
CREATE INDEX IF NOT EXISTS idx_delivery_deadline ON delivery(deadline);
CREATE INDEX IF NOT EXISTS idx_delivery_courier_deadline ON delivery(courier_id, deadline);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_couriers_status;
DROP INDEX IF EXISTS idx_couriers_phone;
DROP INDEX IF EXISTS idx_couriers_status_available;
DROP INDEX IF EXISTS idx_delivery_deadline;
DROP INDEX IF EXISTS idx_delivery_courier_deadline;
-- +goose StatementEnd
