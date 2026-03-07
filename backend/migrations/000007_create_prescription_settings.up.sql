CREATE TABLE IF NOT EXISTS prescription_settings (
    id                      SERIAL PRIMARY KEY,
    doctor_id               INT NOT NULL REFERENCES doctors(id) ON DELETE CASCADE,
    chamber_id              INT NOT NULL REFERENCES chambers(id) ON DELETE CASCADE,
    
    header_left_bangla      TEXT,
    header_right_english    TEXT,
    footer_info_bangla      TEXT,
    footer_info_english     TEXT,
    template_type           VARCHAR(50) NOT NULL DEFAULT 'standard',

    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at              TIMESTAMPTZ,

    UNIQUE(doctor_id, chamber_id)
);

CREATE INDEX idx_prescription_settings_doctor_id ON prescription_settings(doctor_id);
CREATE INDEX idx_prescription_settings_chamber_id ON prescription_settings(chamber_id);
