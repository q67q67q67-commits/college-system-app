-- Уведомления и дополнительные таблицы
CREATE TABLE IF NOT EXISTS notifications (
    id              BIGSERIAL PRIMARY KEY,
    title           TEXT NOT NULL,
    body            TEXT,
    created_by      BIGINT REFERENCES users (id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_notifications_created ON notifications (created_at DESC);

-- Посты директора (профиль как Instagram)
CREATE TABLE IF NOT EXISTS director_posts (
    id              BIGSERIAL PRIMARY KEY,
    director_id     BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    media_url       TEXT,
    media_type      TEXT DEFAULT 'image',
    caption         TEXT,
    is_pinned       BOOLEAN NOT NULL DEFAULT false,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_director_posts_director ON director_posts (director_id);

-- Карта здания (для админа/директора) - одна запись
CREATE TABLE IF NOT EXISTS building_map (
    id              BIGSERIAL PRIMARY KEY,
    content         TEXT,
    image_url       TEXT,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
