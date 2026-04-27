-- 000001_create_todos_table.down.sql

DROP INDEX IF EXISTS idx_todos_created_at;
DROP INDEX IF EXISTS idx_todos_priority;
DROP INDEX IF EXISTS idx_todos_done;
DROP TABLE IF EXISTS todos;
