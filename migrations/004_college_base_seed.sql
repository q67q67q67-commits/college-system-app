-- 004_college_base_seed.sql
-- Укороченный сидинг по данным из college_base.md
-- Базовые группы, студенты и преподаватели (без расписания).
-- Пароль для всех аккаунтов: password123

-- =========================
-- Группы (по кодам college_base.md)
-- =========================
INSERT INTO groups (name, description, academic_year) VALUES
    -- 1 курс
    ('БжБ-131',   'Банк және Басқару, 1 курс',                               '2024-2025'),
    ('МНК-131',   'Менеджмент, 1 курс',                                      '2024-2025'),
    ('МРК-131',   'Маркетинг, 1 курс',                                       '2024-2025'),
    ('БҚ-131',    'Бағдарламалық қамтамасыз ету, 1 курс',                    '2024-2025'),

    -- 2 курс
    ('БI-231',    'Банк ісі, 2 курс',                                        '2023-2024'),
    ('ЕжА-231',   'Есеп және Аудит, 2 курс',                                 '2023-2024'),
    ('МНК-231',   'Менеджмент, 2 курс',                                      '2023-2024'),
    ('МРК-231',   'Маркетинг, 2 курс',                                       '2023-2024'),
    ('ҚҚ-231',    'Құқықтану, 2 курс',                                       '2023-2024'),
    ('БҚ-231',    'Бағдарламалық қамтамасыз ету, 2 курс',                    '2023-2024'),
    ('ЕТжАЖ-231', 'Есептеу техникасы және ақпараттық желілер, 2 курс',      '2023-2024'),

    -- 3 курс
    ('БI-331',    'Банк ісі, 3 курс',                                        '2022-2023'),
    ('ЕжА-331',   'Есеп және Аудит, 3 курс',                                 '2022-2023'),
    ('МНК-331',   'Менеджмент, 3 курс',                                      '2022-2023'),
    ('МРК-331',   'Маркетинг, 3 курс',                                       '2022-2023'),
    ('ҚҚ-331',    'Құқықтану, 3 курс',                                       '2022-2023'),
    ('БҚ-331',    'Бағдарламалық қамтамасыз ету, 3 курс',                    '2022-2023'),
    ('ЕТжАЖ-331', 'Есептеу техникасы және ақпараттық желілер, 3 курс',      '2022-2023');


-- =========================
-- Преподаватели
-- =========================
-- 1 курс (общие предметы)
INSERT INTO users (email, password_hash, role, full_name, language) VALUES
    ('makataeva@college-narxoz.kz',
     crypt('password123', gen_salt('bf')), 'teacher',
     'Мақатаева Жанфия Дәулетқызы', 'kk'),
    ('ermaganbetov@college-narxoz.kz',
     crypt('password123', gen_salt('bf')), 'teacher',
     'Ермаганбетов Кайрат Курманович', 'kk'),
    ('edigeeva@college-narxoz.kz',
     crypt('password123', gen_salt('bf')), 'teacher',
     'Едігеева Мәншүк Нұржауқызы', 'kk');

-- БҚ-331 (профильные преподы по college_base.md)
INSERT INTO users (email, password_hash, role, full_name, language) VALUES
    ('zhanysbai@college-narxoz.kz',
     crypt('password123', gen_salt('bf')), 'teacher',
     'Жанысбай Нұрбақыт Рахатұлы', 'kk'),
    ('zhubanov@college-narxoz.kz',
     crypt('password123', gen_salt('bf')), 'teacher',
     'Жубанов Айбек Нұрланұлы', 'kk');


-- =========================
-- Студенты (укороченный набор: по одному студенту на группу)
-- Email: первая латинская буква (или digraph) имени + фамилия латиницей
-- =========================
INSERT INTO users (email, password_hash, role, full_name, language) VALUES
    -- 1 курс
    -- БжБ-131: Айдар Серік  (фамилия Айдар, имя Серік) -> saidar
    ('saidar@college-narxoz.kz',     crypt('password123', gen_salt('bf')), 'student', 'Айдар Серік', 'kk'),
    -- МНК-131: Абай Адина -> aabay
    ('aabay@college-narxoz.kz',      crypt('password123', gen_salt('bf')), 'student', 'Абай Адина', 'kk'),
    -- МРК-131: Дастанов Улан -> udastanov
    ('udastanov@college-narxoz.kz',  crypt('password123', gen_salt('bf')), 'student', 'Дастанов Улан', 'kk'),
    -- БҚ-131: Жубанов Ахмет -> azhubanov
    ('azhubanov@college-narxoz.kz',  crypt('password123', gen_salt('bf')), 'student', 'Жубанов Ахмет', 'kk'),

    -- 2 курс
    -- БI-231: Қалмырза Шерхан -> shkalmurza
    ('shkalmurza@college-narxoz.kz', crypt('password123', gen_salt('bf')), 'student', 'Қалмырза Шерхан', 'kk'),
    -- ЕжА-231: Алпысбай Кәусар -> kalpysbai
    ('kalpysbai@college-narxoz.kz',  crypt('password123', gen_salt('bf')), 'student', 'Алпысбай Кәусар', 'kk'),
    -- МНК-231: Утепов Ернар -> eutepov
    ('eutepov@college-narxoz.kz',    crypt('password123', gen_salt('bf')), 'student', 'Утепов Ернар', 'kk'),
    -- МРК-231: Олжабай Алима -> aolzhabai
    ('aolzhabai@college-narxoz.kz',  crypt('password123', gen_salt('bf')), 'student', 'Олжабай Алима', 'kk'),
    -- ҚҚ-231: Балтабай Хорлан -> hbaltabai
    ('hbaltabai@college-narxoz.kz',  crypt('password123', gen_salt('bf')), 'student', 'Балтабай Хорлан', 'kk'),
    -- БҚ-231: Адамбай Нуртас -> nadambai
    ('nadambai@college-narxoz.kz',   crypt('password123', gen_salt('bf')), 'student', 'Адамбай Нуртас', 'kk'),
    -- ЕТжАЖ-231: Хамза Айгүл -> akhamza
    ('akhamza@college-narxoz.kz',    crypt('password123', gen_salt('bf')), 'student', 'Хамза Айгүл', 'kk'),

    -- 3 курс
    -- БI-331: Байтығұл Айяжан -> abaitygul
    ('abaitygul@college-narxoz.kz',  crypt('password123', gen_salt('bf')), 'student', 'Байтығұл Айяжан', 'kk'),
    -- ЕжА-331: Боранбаева Мөлдір -> mboranbaeva
    ('mboranbaeva@college-narxoz.kz',crypt('password123', gen_salt('bf')), 'student', 'Боранбаева Мөлдір', 'kk'),
    -- МНК-331: Берікқызы Құралай -> kberikkyzy
    ('kberikkyzy@college-narxoz.kz', crypt('password123', gen_salt('bf')), 'student', 'Берікқызы Құралай', 'kk'),
    -- МРК-331: Ізбасар Бекболат -> bizbasar
    ('bizbasar@college-narxoz.kz',   crypt('password123', gen_salt('bf')), 'student', 'Ізбасар Бекболат', 'kk'),
    -- ҚҚ-331: Арал Арлан -> aarlan
    ('aaral@college-narxoz.kz',      crypt('password123', gen_salt('bf')), 'student', 'Арал Арлан', 'kk'),
    -- БҚ-331: Абдугаппаров Нурамир -> nabdugapparov
    ('nabdugapparov@college-narxoz.kz', crypt('password123', gen_salt('bf')), 'student', 'Абдугаппаров Нурамир', 'kk'),
    -- ЕТжАЖ-331: Бегжанов Ильяс -> ibegzhanov
    ('ibegzhanov@college-narxoz.kz', crypt('password123', gen_salt('bf')), 'student', 'Бегжанов Ильяс', 'kk');


-- =========================
-- Привязка студентов к группам (user_groups)
-- =========================
INSERT INTO user_groups (user_id, group_id, is_curator)
VALUES
    -- 1 курс
    ((SELECT id FROM users WHERE email='saidar@college-narxoz.kz'),
     (SELECT id FROM groups WHERE name='БжБ-131'), false),
    ((SELECT id FROM users WHERE email='aabay@college-narxoz.kz'),
     (SELECT id FROM groups WHERE name='МНК-131'), false),
    ((SELECT id FROM users WHERE email='udastanov@college-narxoz.kz'),
     (SELECT id FROM groups WHERE name='МРК-131'), false),
    ((SELECT id FROM users WHERE email='azhubanov@college-narxoz.kz'),
     (SELECT id FROM groups WHERE name='БҚ-131'), false),

    -- 2 курс
    ((SELECT id FROM users WHERE email='shkalmurza@college-narxoz.kz'),
     (SELECT id FROM groups WHERE name='БI-231'), false),
    ((SELECT id FROM users WHERE email='kalpysbai@college-narxoz.kz'),
     (SELECT id FROM groups WHERE name='ЕжА-231'), false),
    ((SELECT id FROM users WHERE email='eutepov@college-narxoz.kz'),
     (SELECT id FROM groups WHERE name='МНК-231'), false),
    ((SELECT id FROM users WHERE email='aolzhabai@college-narxoz.kz'),
     (SELECT id FROM groups WHERE name='МРК-231'), false),
    ((SELECT id FROM users WHERE email='hbaltabai@college-narxoz.kz'),
     (SELECT id FROM groups WHERE name='ҚҚ-231'), false),
    ((SELECT id FROM users WHERE email='nadambai@college-narxoz.kz'),
     (SELECT id FROM groups WHERE name='БҚ-231'), false),
    ((SELECT id FROM users WHERE email='akhamza@college-narxoz.kz'),
     (SELECT id FROM groups WHERE name='ЕТжАЖ-231'), false),

    -- 3 курс
    ((SELECT id FROM users WHERE email='abaitygul@college-narxoz.kz'),
     (SELECT id FROM groups WHERE name='БI-331'), false),
    ((SELECT id FROM users WHERE email='mboranbaeva@college-narxoz.kz'),
     (SELECT id FROM groups WHERE name='ЕжА-331'), false),
    ((SELECT id FROM users WHERE email='kberikkyzy@college-narxoz.kz'),
     (SELECT id FROM groups WHERE name='МНК-331'), false),
    ((SELECT id FROM users WHERE email='bizbasar@college-narxoz.kz'),
     (SELECT id FROM groups WHERE name='МРК-331'), false),
    ((SELECT id FROM users WHERE email='aaral@college-narxoz.kz'),
     (SELECT id FROM groups WHERE name='ҚҚ-331'), false),
    ((SELECT id FROM users WHERE email='nabdugapparov@college-narxoz.kz'),
     (SELECT id FROM groups WHERE name='БҚ-331'), false),
    ((SELECT id FROM users WHERE email='ibegzhanov@college-narxoz.kz'),
     (SELECT id FROM groups WHERE name='ЕТжАЖ-331'), false);

