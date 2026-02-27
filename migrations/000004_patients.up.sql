-- 000004_create_patients_table.up.sql

CREATE TYPE blood_group_type AS ENUM (
    'A+',
    'A-',
    'B+',
    'B-',
    'AB+',
    'AB-',
    'O+',
    'O-'
);

CREATE TABLE IF NOT EXISTS patients (
    id                  SERIAL PRIMARY KEY,
    doctor_id           INT NOT NULL REFERENCES doctors(id) ON DELETE RESTRICT,

    -- Basic Info
    full_name           VARCHAR(150)        NOT NULL,
    gender              gender_type         NOT NULL,
    age                 SMALLINT            NOT NULL,

    -- Contact
    phone               VARCHAR(20),
    email               VARCHAR(255),

    -- Medical Info
    blood_group         blood_group_type,
    allergies           TEXT[],
    medical_history     TEXT,

    created_at          TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ
);

-- INDEXES
CREATE INDEX idx_patients_doctor_id         ON patients (doctor_id);
CREATE INDEX idx_patients_full_name         ON patients (full_name);
CREATE INDEX idx_patients_phone             ON patients (phone);
CREATE INDEX idx_patients_blood_group       ON patients (blood_group);
CREATE INDEX idx_patients_deleted_at        ON patients (deleted_at)
    WHERE deleted_at