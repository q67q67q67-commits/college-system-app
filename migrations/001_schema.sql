-- Narxoz College (NC) — базовая схема БД
-- Роли: student, teacher, director, admin
-- Лимит хранилища на пользователя: 2 ГБ
-- Анонимность постов только в форуме

-- Расширение для UUID (опционально, можно заменить на bigserial)
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Роли пользователей (enum для типобезопасности)
CREATE TYPE user_role AS ENUM ('student', 'teacher', 'director', 'admin');

-- =============================================================================
-- Пользователи (с ролями и лимитом хранилища 2 ГБ)
-- =============================================================================
CREATE TABLE users (
    id              BIGSERIAL PRIMARY KEY,
    email           TEXT NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL,
    role            user_role NOT NULL DEFAULT 'student',
    full_name       TEXT NOT NULL,
    phone           TEXT,
    avatar_url      TEXT,
    language        TEXT NOT NULL DEFAULT 'ru',
    -- Лимит хранилища: 2 ГБ на участника (в байтах)
    storage_limit_bytes  BIGINT NOT NULL DEFAULT 2147483648,
    storage_used_bytes   BIGINT NOT NULL DEFAULT 0,
    -- Доступ после ухода из заведения (до 1 года)
    left_at         TIMESTAMPTZ,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_users_role ON users (role);
CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_is_active ON users (is_active);

-- =============================================================================
-- Группы (учебные группы)
-- =============================================================================
CREATE TABLE groups (
    id              BIGSERIAL PRIMARY KEY,
    name            TEXT NOT NULL,
    description     TEXT,
    academic_year   TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Связь пользователь — группа (студенты в группах, кураторы и т.д.)
CREATE TABLE user_groups (
    user_id         BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    group_id        BIGINT NOT NULL REFERENCES groups (id) ON DELETE CASCADE,
    is_curator      BOOLEAN NOT NULL DEFAULT false,
    PRIMARY KEY (user_id, group_id)
);

CREATE INDEX idx_user_groups_group ON user_groups (group_id);

-- =============================================================================
-- Расписание (пары, занятия)
-- =============================================================================
CREATE TABLE schedules (
    id              BIGSERIAL PRIMARY KEY,
    group_id        BIGINT NOT NULL REFERENCES groups (id) ON DELETE CASCADE,
    teacher_id      BIGINT NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    subject         TEXT NOT NULL,
    room            TEXT,
    day_of_week     SMALLINT NOT NULL CHECK (day_of_week BETWEEN 1 AND 7),
    start_time      TIME NOT NULL,
    end_time        TIME NOT NULL,
    note            TEXT,
    academic_period TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_schedules_group ON schedules (group_id);
CREATE INDEX idx_schedules_teacher ON schedules (teacher_id);
CREATE INDEX idx_schedules_day ON schedules (day_of_week);

-- Дополнения к парам (ДЗ, уведомления)
CREATE TABLE schedule_attachments (
    id              BIGSERIAL PRIMARY KEY,
    schedule_id     BIGINT NOT NULL REFERENCES schedules (id) ON DELETE CASCADE,
    title           TEXT NOT NULL,
    body            TEXT,
    attachment_type TEXT NOT NULL DEFAULT 'homework', -- homework, notification, material
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- =============================================================================
-- Журнал оценок
-- =============================================================================
CREATE TABLE grades (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    schedule_id     BIGINT NOT NULL REFERENCES schedules (id) ON DELETE CASCADE,
    grade           NUMERIC(5, 2) NOT NULL,
    grade_date      DATE NOT NULL,
    comment         TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_grades_user ON grades (user_id);
CREATE INDEX idx_grades_schedule ON grades (schedule_id);
CREATE INDEX idx_grades_date ON grades (grade_date);

-- =============================================================================
-- Форум (с поддержкой анонимности)
-- =============================================================================
CREATE TABLE forum_posts (
    id              BIGSERIAL PRIMARY KEY,
    parent_id       BIGINT REFERENCES forum_posts (id) ON DELETE CASCADE,
    author_id       BIGINT REFERENCES users (id) ON DELETE SET NULL,
    -- Анонимность: при true автор не отображается, author_id храним для модерации
    is_anonymous    BOOLEAN NOT NULL DEFAULT false,
    title           TEXT,
    body            TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_forum_posts_parent ON forum_posts (parent_id);
CREATE INDEX idx_forum_posts_author ON forum_posts (author_id);
CREATE INDEX idx_forum_posts_created ON forum_posts (created_at DESC);

-- =============================================================================
-- Библиотека (книги, экземпляры, бронирование)
-- =============================================================================
CREATE TABLE library_books (
    id              BIGSERIAL PRIMARY KEY,
    title           TEXT NOT NULL,
    author          TEXT,
    isbn            TEXT,
    description     TEXT,
    cover_url       TEXT,
    has_physical    BOOLEAN NOT NULL DEFAULT true,
    total_copies    INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE library_copies (
    id              BIGSERIAL PRIMARY KEY,
    book_id         BIGINT NOT NULL REFERENCES library_books (id) ON DELETE CASCADE,
    holder_id       BIGINT REFERENCES users (id) ON DELETE SET NULL,
    borrowed_at    TIMESTAMPTZ,
    due_at          DATE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE library_reservations (
    id              BIGSERIAL PRIMARY KEY,
    book_id         BIGINT NOT NULL REFERENCES library_books (id) ON DELETE CASCADE,
    user_id         BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    reserved_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at      TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_library_copies_book ON library_copies (book_id);
CREATE INDEX idx_library_copies_holder ON library_copies (holder_id);

-- =============================================================================
-- События и мероприятия (главная, лента событий)
-- =============================================================================
CREATE TABLE events (
    id              BIGSERIAL PRIMARY KEY,
    title           TEXT NOT NULL,
    description     TEXT,
    image_url       TEXT,
    event_date      TIMESTAMPTZ NOT NULL,
    location        TEXT,
    created_by      BIGINT REFERENCES users (id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_events_date ON events (event_date DESC);

-- Обновление updated_at триггер (шаблон)
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Применяем триггер к таблицам с updated_at
CREATE TRIGGER tr_users_updated_at
    BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER tr_groups_updated_at
    BEFORE UPDATE ON groups FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER tr_schedules_updated_at
    BEFORE UPDATE ON schedules FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER tr_grades_updated_at
    BEFORE UPDATE ON grades FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER tr_forum_posts_updated_at
    BEFORE UPDATE ON forum_posts FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER tr_library_books_updated_at
    BEFORE UPDATE ON library_books FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER tr_events_updated_at
    BEFORE UPDATE ON events FOR EACH ROW EXECUTE FUNCTION set_updated_at();
