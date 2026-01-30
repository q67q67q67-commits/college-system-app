-- Файлы пользователей (Мои файлы), учёт квоты 2 ГБ

CREATE TABLE user_files (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    path            TEXT NOT NULL,
    filename        TEXT NOT NULL,
    size_bytes      BIGINT NOT NULL DEFAULT 0,
    content_type    TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id, path)
);

CREATE INDEX idx_user_files_user ON user_files (user_id);

-- Триггер: при добавлении/удалении/обновлении user_files обновлять users.storage_used_bytes
CREATE OR REPLACE FUNCTION update_user_storage_used()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE users SET storage_used_bytes = storage_used_bytes + NEW.size_bytes WHERE id = NEW.user_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE users SET storage_used_bytes = storage_used_bytes - OLD.size_bytes WHERE id = OLD.user_id;
    ELSIF TG_OP = 'UPDATE' AND OLD.size_bytes IS DISTINCT FROM NEW.size_bytes THEN
        UPDATE users SET storage_used_bytes = storage_used_bytes - OLD.size_bytes + NEW.size_bytes WHERE id = NEW.user_id;
    END IF;
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_user_files_storage
    AFTER INSERT OR UPDATE OR DELETE ON user_files FOR EACH ROW EXECUTE FUNCTION update_user_storage_used();
