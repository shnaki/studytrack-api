-- subjects テーブルを projects にリネーム
ALTER TABLE subjects RENAME TO projects;

-- study_logs の subject_id を project_id にリネーム
ALTER TABLE study_logs RENAME COLUMN subject_id TO project_id;

-- goals の subject_id を project_id にリネーム
ALTER TABLE goals RENAME COLUMN subject_id TO project_id;

-- インデックスのリネーム
ALTER INDEX idx_study_logs_subject RENAME TO idx_study_logs_project;

-- 制約のリネーム
ALTER TABLE projects RENAME CONSTRAINT subjects_pkey TO projects_pkey;
ALTER TABLE projects RENAME CONSTRAINT subjects_user_id_name_key TO projects_user_id_name_key;
ALTER TABLE projects RENAME CONSTRAINT subjects_user_id_fkey TO projects_user_id_fkey;
