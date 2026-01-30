-- Narxoz College (NC) — псевдо-база: группы, пользователи, расписание, оценки, форум, библиотека, события
-- Пароль для всех тестовых аккаунтов: password123

-- Группы (по специальностям с сайта колледжа)
INSERT INTO groups (name, description, academic_year) VALUES
('Банк-1', 'Банковское и страховое дело, 1 курс', '2024-2025'),
('Учёт-1', 'Учёт и аудит, 1 курс', '2024-2025'),
('ПО-1', 'Программное обеспечение, 1 курс', '2024-2025'),
('Менеджмент-1', 'Менеджмент, 1 курс', '2024-2025');

-- Пользователи: директор, преподаватели, студенты (password_hash = bcrypt('password123'))
-- Используем pgcrypto: crypt('password123', gen_salt('bf'))
INSERT INTO users (email, password_hash, role, full_name, phone, language) VALUES
('director@nc.kz', crypt('password123', gen_salt('bf')), 'director', 'Абайдуллаев Мақсат Серікболұлы', '+77068080002', 'ru'),
('admin@nc.kz', crypt('password123', gen_salt('bf')), 'admin', 'Администратор Системы', NULL, 'ru'),
('teacher1@nc.kz', crypt('password123', gen_salt('bf')), 'teacher', 'Иванова Мария Петровна', NULL, 'ru'),
('teacher2@nc.kz', crypt('password123', gen_salt('bf')), 'teacher', 'Сидоров Алексей Николаевич', NULL, 'ru'),
('student1@nc.kz', crypt('password123', gen_salt('bf')), 'student', 'Козлов Артём', NULL, 'ru'),
('student2@nc.kz', crypt('password123', gen_salt('bf')), 'student', 'Нурланова Айгуль', NULL, 'ru'),
('student3@nc.kz', crypt('password123', gen_salt('bf')), 'student', 'Петров Дмитрий', NULL, 'ru'),
('student4@nc.kz', crypt('password123', gen_salt('bf')), 'student', 'Касымова Дана', NULL, 'ru'),
('student5@nc.kz', crypt('password123', gen_salt('bf')), 'student', 'Омаров Нурлан', NULL, 'ru');

-- Связь пользователь — группа (студенты в группах, кураторы)
-- group id: 1 Банк-1, 2 Учёт-1, 3 ПО-1, 4 Менеджмент-1
-- user id: 1 director, 2 admin, 3-4 teachers, 5-9 students
INSERT INTO user_groups (user_id, group_id, is_curator) VALUES
(5, 1, false), (6, 1, false), (7, 2, false), (8, 2, false), (9, 3, false),
(3, 1, true), (3, 2, true), (4, 3, true), (4, 4, true);

-- Расписание (пн=1 .. вс=7)
INSERT INTO schedules (group_id, teacher_id, subject, room, day_of_week, start_time, end_time, academic_period) VALUES
(1, 3, 'Основы банковского дела', '101', 1, '09:00', '10:30', '2024-2025'),
(1, 3, 'Финансовая математика', '101', 1, '10:45', '12:15', '2024-2025'),
(2, 3, 'Бухгалтерский учёт', '102', 2, '09:00', '10:30', '2024-2025'),
(3, 4, 'Программирование на Python', '201', 1, '14:00', '15:30', '2024-2025'),
(3, 4, 'Базы данных', '201', 3, '09:00', '10:30', '2024-2025');

-- Дополнения к парам (ДЗ)
INSERT INTO schedule_attachments (schedule_id, title, body, attachment_type) VALUES
(1, 'ДЗ №1', 'Прочитать главу 1-2 учебника, задачи 1-5', 'homework'),
(4, 'Лаб. работа 1', 'Написать программу «Калькулятор»', 'homework');

-- Оценки (user_id 5,6,7 — студенты; schedule_id 1,2,3,4,5)
INSERT INTO grades (user_id, schedule_id, grade, grade_date, comment) VALUES
(5, 1, 85.5, CURRENT_DATE - 7, NULL),
(5, 2, 90, CURRENT_DATE - 5, NULL),
(6, 1, 78, CURRENT_DATE - 7, NULL),
(7, 3, 92, CURRENT_DATE - 3, NULL),
(8, 3, 88, CURRENT_DATE - 3, NULL),
(9, 4, 95, CURRENT_DATE - 2, 'Отлично');

-- Форум (темы и ответ; один анонимный)
INSERT INTO forum_posts (parent_id, author_id, is_anonymous, title, body) VALUES
(NULL, 5, false, 'Вопрос по сессии', 'Когда расписание экзаменов для Банк-1?'),
(NULL, 7, true, 'Обсуждение мероприятия', 'Кто идёт на встречу с Popeyes?');
INSERT INTO forum_posts (parent_id, author_id, is_anonymous, title, body) VALUES
(1, 3, false, NULL, 'Расписание выложат на следующей неделе.');

-- Библиотека: книги и экземпляры
INSERT INTO library_books (title, author, isbn, has_physical, total_copies) VALUES
('Основы банковского дела', 'Лаврушин О.И.', '978-5-16-012345-6', true, 5),
('Программирование на Go', 'Донован А.', '978-5-97060-415-4', true, 3),
('Бухгалтерский учёт', 'Кондраков Н.П.', '978-5-16-006789-0', true, 4);

INSERT INTO library_copies (book_id, holder_id, borrowed_at, due_at) VALUES
(1, 5, NOW() - INTERVAL '5 days', CURRENT_DATE + 14),
(1, NULL, NULL, NULL),
(2, 9, NOW() - INTERVAL '2 days', CURRENT_DATE + 19),
(3, NULL, NULL, NULL);

INSERT INTO library_reservations (book_id, user_id, reserved_at, expires_at) VALUES
(2, 6, NOW(), NOW() + INTERVAL '3 days');

-- События / новости (главная, лента)
INSERT INTO events (title, description, image_url, event_date, location, created_by) VALUES
('Навигатор в мире финансов — путь к успешной карьере', 'Встреча с экспертами банковской сферы.', NULL, CURRENT_DATE + 14, 'Аудитория 56', 1),
('Popeyes & Narxoz College — партнёрство', 'Презентация программ стажировок.', NULL, CURRENT_DATE + 7, 'Актовый зал', 1),
('Мероприятия по повышению качества образования', 'Итоги семестра и планы на следующий год.', NULL, CURRENT_DATE - 5, NULL, 1);
