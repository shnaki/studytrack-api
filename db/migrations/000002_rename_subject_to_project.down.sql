-- 制約のリネーム（逆）
ALTER TABLE projects RENAME CONSTRAINT projects_user_id_fkey TO subjects_user_id_fkey;
ALTER TABLE projects RENAME CONSTRAINT projects_user_id_name_key TO subjects_user_id_name_key;
ALTER TABLE projects RENAME CONSTRAINT projects_pkey TO subjects_pkey;

-- インデックスのリネーム（逆）
ALTER INDEX idx_study_logs_project RENAME TO idx_study_logs_subject;

-- goals の project_id を subject_id にリネーム
ALTER TABLE goals RENAME COLUMN project_id TO subject_id;

-- study_logs の project_id を subject_id にリネーム
ALTER TABLE study_logs RENAME COLUMN project_id TO subject_id;

-- projects テーブルを subjects にリネーム
ALTER TABLE projects RENAME TO subjects;
