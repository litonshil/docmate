DROP INDEX IF EXISTS idx_doctors_deleted_at;
DROP INDEX IF EXISTS idx_doctors_phone;
DROP INDEX IF EXISTS idx_doctors_full_name;
DROP INDEX IF EXISTS idx_doctors_specialization;
DROP TABLE IF EXISTS doctors;
DROP TYPE IF EXISTS specialization_type;
DROP TYPE IF EXISTS degree_type;