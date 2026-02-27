CREATE TYPE user_role AS ENUM (
    'doctor',
    'admin'
);

CREATE TABLE IF NOT EXISTS users (
                                     id          SERIAL PRIMARY KEY,
                                     username    VARCHAR(255)        NOT NULL UNIQUE,
    email       VARCHAR(255)        NOT NULL UNIQUE,
    password    VARCHAR(255)        NOT NULL,
    role        user_role           NOT NULL DEFAULT 'doctor',

    created_at  TIMESTAMP           NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP           NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMP
    );

-- ============================================================
-- INDEXES
-- ============================================================

-- Role filtering
CREATE INDEX idx_users_role         ON users (role);

-- Soft delete filtering
CREATE INDEX idx_users_deleted_at   ON users (deleted_at)
    WHERE deleted_at IS NULL;