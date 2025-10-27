-- +goose Up
-- создаем расширение для работы с UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- создаем таблицу заказов
CREATE TABLE IF NOT EXISTS orders
(
    uuid            UUID PRIMARY KEY,
    user_id         VARCHAR(255) NOT NULL,
    part_uuids      UUID[] NOT NULL DEFAULT '{}',
    total_price     NUMERIC(12, 2) NOT NULL DEFAULT 0.00,
    transaction_id  VARCHAR(255) DEFAULT '',
    payment_method  VARCHAR(50) DEFAULT '',
    status          VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP
);

-- создаем индексы для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);

-- +goose Down
-- удаляем индексы
DROP INDEX IF EXISTS idx_orders_created_at;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_user_id;

-- удаляем таблицу заказов
DROP TABLE IF EXISTS orders;
