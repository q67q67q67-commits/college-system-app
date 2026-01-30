# Архитектура проекта Narxoz College (NC)

## Структура каталогов

```
NC/
├── cmd/
│   └── api/              # Точка входа API-сервера (Go)
├── internal/             # Внутренние пакеты (domain, usecases, repository, handlers)
├── api/                  # API-спецификации, OpenAPI, контракты
├── design-prototypes/    # Прототипы дизайна для импорта из Figma
├── migrations/           # SQL-миграции (схема и данные)
├── docker-compose.yml    # Локальное окружение (PostgreSQL)
├── go.mod
└── ARCHITECTURE.md
```

## Стек

- **Frontend:** Flutter (мобильная + веб, адаптация под ПК и мобильные)
- **Backend:** Go
- **БД:** PostgreSQL (любой SQL по ТЗ)

## Роли

| Роль       | Описание |
|-----------|----------|
| `student` | Студент  |
| `teacher` | Преподаватель |
| `director`| Директор (профиль как «Twitter» для всех) |
| `admin`   | Админ (без страницы профиля, админ-панель) |

## База данных

- **Пользователи:** роли, лимит хранилища **2 ГБ** на участника (`storage_limit_bytes`, `storage_used_bytes`).
- **Форум:** поддержка анонимности постов (`forum_posts.is_anonymous`); автор хранится для модерации.
- **Таблицы:** `users`, `groups`, `schedules`, `grades`, `forum_posts`, `library_books` (+ копии и брони), `events`.
- Миграции применяются при первом запуске контейнера PostgreSQL из `migrations/`.

## Дизайн

Макеты и экспорты из Figma размещаются в `design-prototypes/` и подключаются по мере готовности.
