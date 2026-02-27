DROP INDEX IF EXISTS idx_patients_deleted_at;
DROP INDEX IF EXISTS idx_patients_blood_group;
DROP INDEX IF EXISTS idx_patients_phone;
DROP INDEX IF EXISTS idx_patients_full_name;
DROP INDEX IF EXISTS idx_patients_doctor_id;
DROP TABLE IF EXISTS patients;
DROP TYPE IF EXISTS blood_group_type;
DROP TYPE IF EXISTS gender_type;