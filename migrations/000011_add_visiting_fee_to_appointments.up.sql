-- Add visiting_fee and is_fee_collected columns to appointments table
ALTER TABLE appointments ADD COLUMN IF NOT EXISTS visiting_fee DECIMAL(10, 2) DEFAULT 0;
ALTER TABLE appointments ADD COLUMN IF NOT EXISTS is_fee_collected BOOLEAN DEFAULT FALSE;
