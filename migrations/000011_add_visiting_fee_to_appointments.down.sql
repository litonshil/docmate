-- Remove visiting_fee and is_fee_collected columns from appointments table
ALTER TABLE appointments DROP COLUMN IF EXISTS visiting_fee;
ALTER TABLE appointments DROP COLUMN IF EXISTS is_fee_collected;
