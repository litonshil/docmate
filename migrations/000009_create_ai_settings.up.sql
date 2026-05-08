CREATE TABLE IF NOT EXISTS ai_settings (
    id SERIAL PRIMARY KEY,
    doctor_id INTEGER UNIQUE NOT NULL REFERENCES doctors(id) ON DELETE CASCADE,
    is_ai_enabled BOOLEAN DEFAULT FALSE,
    allow_global_api BOOLEAN DEFAULT FALSE,
    provider VARCHAR(50) DEFAULT 'gemini',
    individual_api_key TEXT,
    use_individual_key BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_ai_settings_doctor_id ON ai_settings(doctor_id);
