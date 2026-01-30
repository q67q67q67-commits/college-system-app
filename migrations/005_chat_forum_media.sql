-- Чат: сообщения с возможностью удаления
CREATE TABLE IF NOT EXISTS chat_messages (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    body            TEXT NOT NULL,
    media_url       TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_chat_messages_created ON chat_messages (created_at DESC);

-- Форум: медиа в постах
ALTER TABLE forum_posts ADD COLUMN IF NOT EXISTS media_url TEXT;
