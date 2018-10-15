
ALTER TABLE task
ADD COLUMN task_id text,
ADD COLUMN command text,
ADD COLUMN memory integer,
ADD COLUMN user_tags jsonb;

CREATE INDEX IF NOT EXISTS ix_task_task_id ON task(task_id);
CREATE INDEX IF NOT EXISTS ix_task_user_tags ON task USING GIN (user_tags);
