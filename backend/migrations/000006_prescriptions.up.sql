CREATE TABLE IF NOT EXISTS prescriptions (
    id                  SERIAL PRIMARY KEY,
    doctor_id           INT NOT NULL REFERENCES doctors(id) ON DELETE RESTRICT,
    patient_id          INT NOT NULL REFERENCES patients(id) ON DELETE RESTRICT,
    chamber_id          INT NOT NULL REFERENCES chambers(id) ON DELETE RESTRICT,

    -- Vitals
    -- e.g.
    -- {
    --   "weight_kg":        70,
    --   "height_cm":        175,
    --   "blood_pressure":   "120/80",
    --   "temperature_f":    98.6,
    --   "pulse_bpm":        72,
    --   "spo2_percent":     98
    -- }
    vitals              JSONB               NOT NULL DEFAULT '{}',

    -- Chief Complaints
    -- e.g. '{"Headache", "Fever", "Nausea"}'
    chief_complaints    TEXT[]              NOT NULL DEFAULT '{}',

    -- Diagnosis
    -- e.g. '{"Hypertension", "Type 2 Diabetes"}'
    diagnosis           TEXT[]              NOT NULL DEFAULT '{}',

    -- Medications
    -- e.g.
    -- [
    --   {
    --     "medicine_id":    12,               -- optional, ref to medicines table
    --     "medicine_name":  "Napa",
    --     "generic_name":   "Paracetamol",
    --     "form":           "tablet",         -- tablet, capsule, syrup, injection etc.
    --     "strength":       "500mg",
    --     "dosage":         "1 tablet",
    --     "frequency":      "twice daily",
    --     "timing":         "after meal",
    --     "duration":       "5 days",
    --     "instructions":   "with warm water",
    --     "sort_order":     1
    --   }
    -- ]
    medications         JSONB               NOT NULL DEFAULT '[]',

    -- Investigations
    -- e.g. '{"CBC", "Blood Sugar (F)", "Chest X-Ray"}'
    investigations      TEXT[]              NOT NULL DEFAULT '{}',

    -- Advice & Follow Up
    advice              TEXT,
    follow_up_date      DATE,
    file_url            VARCHAR(255) NULL,

    created_at          TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ
);

-- INDEXES
CREATE INDEX idx_prescriptions_doctor_id        ON prescriptions (doctor_id);
CREATE INDEX idx_prescriptions_patient_id       ON prescriptions (patient_id);
CREATE INDEX idx_prescriptions_chamber_id       ON prescriptions (chamber_id);
CREATE INDEX idx_prescriptions_follow_up_date   ON prescriptions (follow_up_date);
CREATE INDEX idx_prescriptions_chief_complaints ON prescriptions USING GIN (chief_complaints);
CREATE INDEX idx_prescriptions_diagnosis        ON prescriptions USING GIN (diagnosis);
CREATE INDEX idx_prescriptions_medications      ON prescriptions USING GIN (medications);
CREATE INDEX idx_prescriptions_investigations   ON prescriptions USING GIN (investigations);
CREATE INDEX idx_prescriptions_deleted_at       ON prescriptions (deleted_at);