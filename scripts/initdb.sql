-- ============================================================
-- DATABASE & USER SETUP
-- ============================================================

-- Create the database
CREATE DATABASE docmate;

-- Create the user with password
CREATE USER docmate_user WITH PASSWORD '12345678';

-- Grant all privileges on the database
GRANT ALL PRIVILEGES ON DATABASE docmate TO docmate_user;

-- Grant schema privileges
GRANT ALL PRIVILEGES ON SCHEMA public TO docmate_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO docmate_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO docmate_user;

-- Connect to the database before running the rest
\c docmate;

-- ============================================================
-- ENUM TYPES
-- ============================================================

CREATE TYPE degree_type AS ENUM (
    'MBBS',
    'BDS',
    'MD',
    'MS',
    'FCPS',
    'MRCP',
    'FRCS',
    'MPH',
    'PhD',
    'DPM',
    'other'
);

CREATE TYPE gender_type AS ENUM (
    'male',
    'female',
    'other'
);


-- ============================================================
-- 1. DOCTORS  (Authentication & Profile)
-- ============================================================

CREATE TABLE doctors (
    id              SERIAL PRIMARY KEY,
    email           VARCHAR(255) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,

    full_name       VARCHAR(150) NOT NULL,
    degree          degree_type[],              -- array: e.g. '{MBBS, FCPS}'
    specialization  VARCHAR(150),
    phone           VARCHAR(20),
    bio             VARCHAR(500),
    signature_url   VARCHAR(500),

    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    email_verified  BOOLEAN NOT NULL DEFAULT FALSE,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ                             -- soft delete
);


-- ============================================================
-- 2. PASSWORD RESET TOKENS
-- ============================================================

CREATE TABLE password_reset_tokens (
    id          SERIAL PRIMARY KEY,
    doctor_id   INT NOT NULL REFERENCES doctors(id) ON DELETE CASCADE,
    token_hash  VARCHAR(255) NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    used_at     TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- ============================================================
-- 3. CHAMBERS / HOSPITALS
-- ============================================================

CREATE TABLE chambers (
    id              SERIAL PRIMARY KEY,
    doctor_id       INT NOT NULL REFERENCES doctors(id) ON DELETE CASCADE,

    name            VARCHAR(200) NOT NULL,
    address_line1   VARCHAR(255),
    address_line2   VARCHAR(255),
    city            VARCHAR(100),
    district        VARCHAR(100),
    postal_code     VARCHAR(20),
    phone           VARCHAR(20),
    timing_notes    VARCHAR(300),
    is_primary      BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- ============================================================
-- 4. DOCTOR SETTINGS
-- ============================================================

CREATE TABLE doctor_settings (
    id                          SERIAL PRIMARY KEY,
    doctor_id                   INT UNIQUE NOT NULL REFERENCES doctors(id) ON DELETE CASCADE,

    default_prescription_footer VARCHAR(500),
    default_advice              TEXT,               -- can be long multi-line text
    timezone                    VARCHAR(60)  NOT NULL DEFAULT 'Asia/Dhaka',
    date_format                 VARCHAR(20)  NOT NULL DEFAULT 'DD/MM/YYYY',

    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- ============================================================
-- 5. PATIENTS
-- ============================================================

CREATE TABLE patieid              SERIAL PRIMARYdoctor_id       INT NOT NULL REFERENCES doctors(id) ON DELETE CASCfull_name       VARCHAR(150) NOT date_of_birth   age_years       SMALgender          gender_phone           VARCHARemail           VARCHAR(address         VARCHAR(blood_group     VARCHAR(5),                 -- A+, B-, O+,notes           TEXT,                       -- free-form clinical nis_active       BOOLEAN NOT NULL DEFAULT Tcreated_at      TIMESTAMPTZ NOT NULL DEFAULT Nupdated_at      TIMESTAMPTZ NOT NULL DEFAULT Ndeleted_at      TIMESTAMPTZ                         -- soft delete
);

CREATE INDEX idx_patients_doctor ON patients(doctor_id);
CREATE INDEX idx_patients_phone  ON patients(phone);


-- ============================================================
-- 6. PRESCRIPTIONS
-- ============================================================

-- medicines JSONB structure (array of objects):
-- [
--   {
--     "medicine_id"  : 12,             -- optional, ref to medicines table
--     "medicine_name": "Napa",
--     "generic_name" : "Paracetamol",
--     "form"         : "tablet",
--     "strength"     : "500mg",
--     "dosage"       : "1 tablet",
--     "frequency"    : "twice daily",
--     "timing"       : "after meal",
--     "duration"     : "5 days",
--     "instructions" : "with warm water",
--     "sort_order"   : 1
--   }
-- ]

CREATE TABLE prescriptions (
    id                  SERIAL PRIMARY KEY,
    patient_id          INT NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    doctor_id           INT NOT NULL REFERENCES doctors(id),
    chamber_id          INT REFERENCES chambers(id),

    next_visit_d
    -- Clinical
    chief_complaints    TEXT,
    diagnosis           TEXT,
    advice              TEXT,
    notes               TEXT,

    -- Medicines as JSON array
    medicines           JSONB NOT NULL DEFAULT '[]',

    -- Vitals
    weight_kg           NUMERIC(5,2),
    bp_systolic         SMALLINT,
    bp_diastolic        SMALLINT,
    temperature_c       NUMERIC(4,1),
    pulse_bpm           SMALLINT,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_prescriptions_patient   ON prescriptions(patient_id);
CREATE INDEX idx_prescriptions_doctor    ON prescriptions(doctor_id);


-- ============================================================
-- 7. MEDICINE MASTER LIST  (Global formulary for autocomplete)
-- ============================================================

CREATE TABLE medicid              SERIAL PRIMARYbrand_name      VARCHAR(200) NOTgeneric_name    VARCHARform            VARCHAR(50),            -- tablet, syrup, capsule, injestrength        VARCHAR(50),            -- e.g. "500mg", "5mmanufacturer    VARCHAR(is_active       BOOLEAN NOT NULL DEFAULTcreated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- ============================================================
-- 9. GENERATED PDFs
-- ============================================================

CREATE TABLE prescription_pdfs (
    id                  SERIAL PRIMARY KEY,
    prescription_id     INT NOT NULL REFERENCES prescriptions(id) ON DELETE CASCADE,
    doctor_id           INT NOT NULL REFERENCES doctors(id),

    file_url        VARCHAR(500) NOT NULL,
    file_size_bytes INT,
    generated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    downloaded_at   TIMESTAMPTZ
);

CREATE INDEX idx_pdfs_prescription ON prescription_pdfs(prescription_id);