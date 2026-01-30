-- Мои заметки (приватные, только владелец)
CREATE TABLE IF NOT EXISTS user_notes (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    title           TEXT NOT NULL DEFAULT '',
    body            TEXT NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_user_notes_user ON user_notes (user_id);

-- Комментарии к заметкам (дополнения)
CREATE TABLE IF NOT EXISTS user_note_comments (
    id              BIGSERIAL PRIMARY KEY,
    note_id         BIGINT NOT NULL REFERENCES user_notes (id) ON DELETE CASCADE,
    user_id         BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    body            TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_user_note_comments_note ON user_note_comments (note_id);

-- Посты директора: добавляем title, body (как в форуме), комментарии
ALTER TABLE director_posts ADD COLUMN IF NOT EXISTS title TEXT DEFAULT '';
ALTER TABLE director_posts ADD COLUMN IF NOT EXISTS body TEXT DEFAULT '';

CREATE TABLE IF NOT EXISTS director_post_comments (
    id              BIGSERIAL PRIMARY KEY,
    post_id         BIGINT NOT NULL REFERENCES director_posts (id) ON DELETE CASCADE,
    user_id         BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    body            TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_director_post_comments_post ON director_post_comments (post_id);
