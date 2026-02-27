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

CREATE TYPE specialization_type AS ENUM (
    -- Medicine
    'General Practice',
    'Cardiology',
    'Neurology',
    'Gastroenterology',
    'Pulmonology',
    'Nephrology',
    'Endocrinology',
    'Rheumatology',
    'Hematology',
    'Oncology',
    'Infectious Disease',

    -- Surgery
    'General Surgery',
    'Neurosurgery',
    'Cardiothoracic Surgery',
    'Orthopedic Surgery',
    'Plastic Surgery',
    'Vascular Surgery',
    'Urology',

    -- Other Specialties
    'Pediatrics',
    'Gynecology',
    'Obstetrics',
    'Psychiatry',
    'Dermatology',
    'Ophthalmology',
    'ENT',
    'Radiology',
    'Anesthesiology',
    'Pathology',
    'Emergency Medicine',
    'ICU / Critical Care',
    'Dentistry',
    'Orthodontics',

    'other'
);

-- ============================================================
-- 1. DOCTORS  (Profile)
-- ============================================================

CREATE TABLE doctors (
    id              SERIAL PRIMARY KEY,
    user_id         INT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    email           VARCHAR(255) UNIQUE NOT NULL,
    full_name       VARCHAR(150) NOT NULL,
    degree          degree_type[],              -- e.g. '{MBBS, FCPS}'
    specialization  specialization_type[],      -- e.g. '{Cardiology, "General Practice"}'
    phone           VARCHAR(20),
    bio             VARCHAR(500),
    signature_url   VARCHAR(500),

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ                             -- soft delete
);

-- ============================================================
-- INDEXES
-- ============================================================

CREATE INDEX idx_doctors_specialization ON doctors USING GIN (specialization);
CREATE INDEX idx_doctors_full_name      ON doctors (full_name);
CREATE INDEX idx_doctors_phone          ON doctors (phone);
CREATE INDEX idx_doctors_deleted_at     ON doctors (deleted_at)
    WHERE deleted_at IS NULL;