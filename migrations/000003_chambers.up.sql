CREATE TABLE IF NOT EXISTS chambers (
    id                  SERIAL PRIMARY KEY,
    doctor_id           INT NOT NULL REFERENCES doctors(id) ON DELETE CASCADE,

    -- Location
    name                VARCHAR(255) NOT NULL,          -- e.g. 'Dhaka Medical Chamber'
    address             VARCHAR(500) NOT NULL,
    area                VARCHAR(150),                   -- e.g. 'Dhanmondi'
    city                VARCHAR(100) NOT NULL,
    country             VARCHAR(100) NOT NULL DEFAULT 'Bangladesh',

    -- Contact
    phone               VARCHAR(20),
    email               VARCHAR(255),

    -- Consultation Fee
    fee                 NUMERIC(10, 2) NOT NULL,
    follow_up_fee       NUMERIC(10, 2),

    -- Visiting Schedule
    -- e.g.
    -- [
    --   {
    --     "day": "saturday",
    --     "slots": [
    --       { "start_time": "09:00", "end_time": "13:00" },
    --       { "start_time": "17:00", "end_time": "21:00" }
    --     ]
    --   },
    --   {
    --     "day": "monday",
    --     "slots": [
    --       { "start_time": "10:00", "end_time": "14:00" }
    --     ]
    --   }
    -- ]
    visiting_hours      JSONB NOT NULL DEFAULT '[]',
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ
    );

-- INDEXES
CREATE INDEX idx_chambers_doctor_id      ON chambers (doctor_id);
CREATE INDEX idx_chambers_city           ON chambers (city);
CREATE INDEX idx_chambers_area           ON chambers (area);
CREATE INDEX idx_chambers_visiting_hours ON chambers USING GIN (visiting_hours);
CREATE INDEX idx_chambers_deleted_at     ON chambers (deleted_at)
    WHERE deleted_at IS NULL;