CREATE TYPE medicine_form_type AS ENUM (
    'tablet',
    'capsule',
    'syrup',
    'suspension',
    'injection',
    'inhaler',
    'drops',
    'cream',
    'ointment',
    'gel',
    'patch',
    'suppository',
    'powder',
    'sachet',
    'other'
);

CREATE TABLE IF NOT EXISTS medicines (
    id                  SERIAL PRIMARY KEY,
    created_by          INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,

    -- Medicine Info
    brand_name          VARCHAR(255)            NOT NULL,
    generic_name        VARCHAR(255)            NOT NULL,
    form                medicine_form_type      NOT NULL,
    strength            VARCHAR(100),                                   -- e.g. '500mg', '250mg/5ml'
    manufacturer        VARCHAR(255),
    description         TEXT,

    is_active           BOOLEAN                 NOT NULL DEFAULT TRUE,

    created_at          TIMESTAMPTZ             NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ             NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ
);

-- INDEXES
CREATE INDEX idx_medicines_created_by           ON medicines (created_by);
CREATE INDEX idx_medicines_brand_name           ON medicines (brand_name);
CREATE INDEX idx_medicines_generic_name         ON medicines (generic_name);
CREATE INDEX idx_medicines_form                 ON medicines (form);
CREATE INDEX idx_medicines_deleted_at           ON medicines (deleted_at)
    WHERE deleted_at IS NULL;

-- Fuzzy search on brand and generic name
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_medicines_brand_name_trgm      ON medicines USING GIN (brand_name gin_trgm_ops);
CREATE INDEX idx_medicines_generic_name_trgm    ON medicines USING GIN (generic_name gin_trgm_ops);