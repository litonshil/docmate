-- Drop legacy tables
DROP TABLE IF EXISTS ai_settings CASCADE;
DROP TABLE IF EXISTS global_settings CASCADE;

-- Create ai_providers table (Admin controlled settings/keys)
CREATE TABLE IF NOT EXISTS ai_providers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    api_key VARCHAR(255),
    model VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Pre-populate default providers
INSERT INTO ai_providers (name, slug, api_key, model, is_active) VALUES
('Google Gemini', 'gemini', '', 'gemini-1.5-flash', true),
('ChatGPT', 'chatgpt', '', 'gpt-4o-mini', true)
ON CONFLICT (slug) DO NOTHING;

-- Create doctor_ai_settings table (Doctor specific settings)
CREATE TABLE IF NOT EXISTS doctor_ai_settings (
    id SERIAL PRIMARY KEY,
    doctor_id INT NOT NULL REFERENCES doctors(id) ON DELETE CASCADE,
    ai_providers_id INT NOT NULL REFERENCES ai_providers(id) ON DELETE CASCADE,
    individual_api_key VARCHAR(255),
    is_active BOOLEAN DEFAULT FALSE,
    is_ai_enabled BOOLEAN DEFAULT FALSE,
    allow_global_api BOOLEAN DEFAULT FALSE,
    use_individual_key BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(doctor_id, ai_providers_id)
);
