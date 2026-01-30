# API Narxoz College (NC) — список эндпоинтов и примеры JSON

Базовый URL: `http://localhost:8080`  
Авторизация: заголовок `Authorization: Bearer <JWT>` (кроме POST /auth/login).

---

## Публичные

### POST /auth/login

Вход. Возвращает JWT и данные пользователя.

**Request:**
```json
{
  "email": "student1@nc.kz",
  "password": "password123"
}
```

**Response 200:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user_id": 5,
  "role": "student",
  "email": "student1@nc.kz",
  "full_name": "Козлов Артём",
  "expires_at": "2026-01-31T12:00:00Z"
}
```

**Response 401:** `{"error":"invalid email or password"}`

---

## С авторизацией (любая роль)

### GET /api/schedule

Расписание текущего пользователя (студент — по группе, преподаватель — по teacher_id).

**Response 200:** массив занятий
```json
[
  {
    "id": 1,
    "group_id": 1,
    "teacher_id": 3,
    "subject": "Основы банковского дела",
    "room": "101",
    "day_of_week": 1,
    "start_time": "09:00:00",
    "end_time": "10:30:00",
    "note": "",
    "academic_period": "2024-2025"
  }
]
```

### GET /api/schedule/{id}/attachments

Дополнения к паре (ДЗ, уведомления). `{id}` — ID занятия.

**Response 200:** массив
```json
[
  {
    "id": 1,
    "schedule_id": 1,
    "title": "ДЗ №1",
    "body": "Прочитать главу 1-2...",
    "attachment_type": "homework",
    "created_at": "2026-01-30T10:00:00Z"
  }
]
```

### GET /api/grades

Оценки текущего пользователя (студент).

**Response 200:** массив
```json
[
  {
    "id": 1,
    "user_id": 5,
    "schedule_id": 1,
    "grade": 85.5,
    "grade_date": "2026-01-23",
    "comment": "",
    "subject": "Основы банковского дела"
  }
]
```

### GET /api/grades/gpa

Средний балл студента.

**Response 200:** `{"gpa": 87.75}`

### GET /api/grades/transcript

Транскрипт (тот же формат, что GET /api/grades).

### GET /api/groups

Список групп с количеством студентов (для преподавателя/админа).

**Response 200:** массив
```json
[
  {
    "id": 1,
    "name": "Банк-1",
    "description": "Банковское и страховое дело, 1 курс",
    "academic_year": "2024-2025",
    "student_count": 2
  }
]
```

### GET /api/groups/{id}/students

Студенты группы. `{id}` — ID группы.

**Response 200:** массив
```json
[
  {
    "user_id": 5,
    "full_name": "Козлов Артём",
    "email": "student1@nc.kz",
    "is_curator": false
  }
]
```

### GET /api/library/books?q=...&limit=20

Поиск книг. `q` — подстрока по названию/автору.

**Response 200:** массив
```json
[
  {
    "id": 1,
    "title": "Основы банковского дела",
    "author": "Лаврушин О.И.",
    "isbn": "978-5-16-012345-6",
    "description": "",
    "has_physical": true,
    "total_copies": 5,
    "available": 4
  }
]
```

### GET /api/library/books/{id}/copies

Экземпляры книги с владельцами. `{id}` — ID книги.

**Response 200:** массив
```json
[
  {
    "id": 1,
    "book_id": 1,
    "holder_id": 5,
    "holder_name": "Козлов Артём",
    "borrowed_at": "2026-01-25T10:00:00Z",
    "due_at": "2026-02-08"
  }
]
```

### POST /api/library/reservations

Бронирование книги.

**Request:**
```json
{
  "book_id": 2,
  "expires_at": "2026-02-02T12:00:00Z"
}
```
`expires_at` опционален (по умолчанию +72ч).

**Response 201:** `{"id": 1}`

### GET /api/forum/posts?parent_id=...&limit=20&offset=0

Список постов: корневые темы (без `parent_id`) или ответы к теме (`parent_id` — ID темы).

**Response 200:** массив (при `is_anonymous: true` поля author_id/author_name не отдаются)
```json
[
  {
    "id": 1,
    "parent_id": null,
    "author_id": 5,
    "author_name": "Козлов Артём",
    "is_anonymous": false,
    "title": "Вопрос по сессии",
    "body": "Когда расписание экзаменов для Банк-1?",
    "created_at": "2026-01-30T10:00:00Z"
  }
]
```

### POST /api/forum/posts

Создать тему или ответ (с опцией анонимно).

**Request:**
```json
{
  "parent_id": null,
  "is_anonymous": false,
  "title": "Тема обсуждения",
  "body": "Текст поста"
}
```
Для ответа: `parent_id` — ID темы. При анонимном посте фронт не отображает автора.

**Response 201:** `{"id": 3}`

### GET /api/events?from=...&to=...&q=...&limit=20&offset=0

Список событий/новостей. Фильтры: `from`, `to` (даты), `q` (поиск по title/description).

**Response 200:** массив
```json
[
  {
    "id": 1,
    "title": "Навигатор в мире финансов",
    "description": "Встреча с экспертами...",
    "image_url": "",
    "event_date": "2026-02-13T00:00:00Z",
    "location": "Аудитория 56",
    "created_by": 1,
    "created_at": "2026-01-30T10:00:00Z"
  }
]
```

### GET /api/events/{id}

Одно событие по ID.

**Response 200:** объект события (как элемент массива выше).  
**Response 404:** `{"error":"not found"}`

### GET /api/files

Список файлов пользователя и квота хранилища.

**Response 200:**
```json
{
  "files": [
    {
      "id": 1,
      "path": "5/document.pdf",
      "filename": "document.pdf",
      "size_bytes": 1024,
      "content_type": "application/pdf"
    }
  ],
  "storage_limit": 2147483648,
  "storage_used": 1024
}
```

### POST /api/files/upload

Загрузка файла. `Content-Type: multipart/form-data`, поле `file`.

**Response 201:** `{"id": 1, "path": "5/document.pdf", "size": 1024}`  
**Response 403:** `{"error":"storage quota exceeded (2 GB limit)"}`

### DELETE /api/files/{id}

Удалить свой файл. `{id}` — ID записи в user_files.

**Response 200:** `{"status":"ok"}`

### GET /api/chat/ws

WebSocket для чата. После апгрейда клиент отправляет/получает JSON, например:  
`{"type":"message","body":"текст"}` — рассылается всем подключённым.  
Требуется заголовок `Authorization: Bearer <JWT>` при первом HTTP-запросе апгрейда (или query `?token=...`).

### GET /api/director

Публичный профиль директора (без авторизации). Для страницы «Профиль директора».

**Response 200:**
```json
{
  "id": 1,
  "full_name": "Абайдуллаев Мақсат Серікболұлы",
  "avatar_url": "",
  "phone": "+77068080002",
  "email": "director@nc.kz"
}
```
**Response 404:** `{"error":"not found"}` — если нет пользователя с ролью director.

### PUT /api/profile

Обновить профиль текущего пользователя: пароль и/или телефон. Только авторизованные. Передавать только те поля, которые нужно изменить.

**Request:**
```json
{
  "password": "новый_пароль",
  "phone": "+7 700 123 45 67"
}
```
Оба поля опциональны. Если поле не передано — не обновляется.

**Response 200:** `{"status":"ok"}`

---

## Только преподаватель (teacher, director, admin)

### POST /api/grades

Добавить оценку студенту.

**Request:**
```json
{
  "user_id": 5,
  "schedule_id": 1,
  "grade": 90,
  "grade_date": "2026-01-30",
  "comment": "Хорошая работа"
}
```

**Response 201:** `{"id": 7}`

### PUT /api/grades/{id}

Изменить оценку.

**Request:** `{"grade": 95, "comment": "Исправлено"}`  
**Response 200:** `{"status":"ok"}`

### DELETE /api/grades/{id}

Удалить оценку.

**Response 200:** `{"status":"ok"}`

### POST /api/schedule/{id}/attachments

Добавить ДЗ/уведомление к паре. `{id}` — ID занятия (schedule_id).

**Request:**
```json
{
  "title": "ДЗ №2",
  "body": "Задачи 6-10",
  "attachment_type": "homework"
}
```
`attachment_type`: `homework`, `notification`, `material`.

**Response 201:** `{"id": 3}`

### PUT /api/attachments/{id}

Изменить ДЗ/уведомление.

**Request:** `{"title": "ДЗ №2 (обновлено)", "body": "Задачи 6-12"}`  
**Response 200:** `{"status":"ok"}`

### DELETE /api/attachments/{id}

Удалить ДЗ/уведомление.

**Response 200:** `{"status":"ok"}`

---

## Только админ и директор (admin, director)

### POST /api/events

Создать событие/новость.

**Request:**
```json
{
  "title": "Новое мероприятие",
  "description": "Описание",
  "image_url": "",
  "event_date": "2026-02-20T14:00:00Z",
  "location": "Актовый зал"
}
```

**Response 201:** `{"id": 4}`

### PUT /api/events/{id}

Изменить событие.

**Request:** те же поля, что в POST.  
**Response 200:** `{"status":"ok"}`

### DELETE /api/events/{id}

Удалить событие.

**Response 200:** `{"status":"ok"}`

---

## Коды ответов и ошибки

- **200** — успех (GET, PUT, DELETE).
- **201** — создано (POST).
- **400** — неверный запрос (тело/параметры): `{"error":"..."}`.
- **401** — нет/неверный JWT: `{"error":"missing Authorization header"}` и т.п.
- **403** — запрещено (роль или квота): `{"error":"forbidden: role not allowed"}` или `{"error":"storage quota exceeded (2 GB limit)"}`.
- **404** — не найдено: `{"error":"not found"}`.
- **500** — внутренняя ошибка: `{"error":"internal error"}`.

Тестовые аккаунты (после сидинга):  
- director@nc.kz, admin@nc.kz, teacher1@nc.kz, teacher2@nc.kz, student1@nc.kz … student5@nc.kz  
Пароль у всех: `password123`.
