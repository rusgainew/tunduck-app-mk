-- +migrate Up
-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'blocked')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP NULL
);

-- Индексы для оптимизации поиска
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);

-- Комментарии для документации
COMMENT ON TABLE users IS 'Таблица пользователей auth-service';
COMMENT ON COLUMN users.id IS 'Уникальный идентификатор пользователя';
COMMENT ON COLUMN users.email IS 'Email пользователя (уникальный)';
COMMENT ON COLUMN users.name IS 'Имя пользователя';
COMMENT ON COLUMN users.password IS 'Хеш пароля (bcrypt)';
COMMENT ON COLUMN users.status IS 'Статус: active, inactive, blocked';
COMMENT ON COLUMN users.last_login IS 'Время последнего входа';

-- +migrate Down
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_status;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
