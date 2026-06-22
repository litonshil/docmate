-- Alter the user_role enum
ALTER TYPE user_role ADD VALUE IF NOT EXISTS 'assistant';

-- Create assistants table
CREATE TABLE IF NOT EXISTS assistants (
    id          SERIAL PRIMARY KEY,
    user_id     INT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    doctor_id   INT NOT NULL REFERENCES doctors(id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    phone       VARCHAR(50) NOT NULL,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create assistant_chambers join table
CREATE TABLE IF NOT EXISTS assistant_chambers (
    assistant_id INT NOT NULL REFERENCES assistants(id) ON DELETE CASCADE,
    chamber_id   INT NOT NULL REFERENCES chambers(id) ON DELETE CASCADE,
    PRIMARY KEY (assistant_id, chamber_id)
);

-- Indexing
CREATE INDEX IF NOT EXISTS idx_assistants_doctor_id ON assistants(doctor_id);
CREATE INDEX IF NOT EXISTS idx_assistant_chambers_chamber_id ON assistant_chambers(chamber_id);
