-- +goose Up
-- +goose StatementBegin
CREATE TABLE delivery (
    id                  BIGSERIAL PRIMARY KEY,
    courier_id          BIGINT NOT NULL REFERENCES couriers(id),
    order_id            VARCHAR(255) NOT NULL UNIQUE,
    assigned_at         TIMESTAMP NOT NULL DEFAULT NOW(),
    deadline            TIMESTAMP NOT NULL
);

CREATE INDEX idx_delivery_courier_id ON delivery(courier_id);
CREATE INDEX idx_delivery_order_id ON delivery(order_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE delivery;
-- +goose StatementEnd