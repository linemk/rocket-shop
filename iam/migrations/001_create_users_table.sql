-- +goose Up
-- создаем расширение для работы с UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- создаем таблицу пользователей
CREATE TABLE IF NOT EXISTS users
(
    user_uuid            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    login                VARCHAR(255) NOT NULL UNIQUE,
    password_hash        VARCHAR(255) NOT NULL,
    email                VARCHAR(255) NOT NULL,
    notification_methods JSONB DEFAULT '[]'::jsonb,
    created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP
);

-- создаем индексы для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_users_login ON users(login);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- +goose Down
-- удаляем индексы
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_login;

-- удаляем таблицу пользователей
DROP TABLE IF EXISTS users;
