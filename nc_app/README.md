# Narxoz College — Flutter (мобильная + веб)

Приложение для студентов и сотрудников колледжа. Сборка для **Web** и **Mobile** (Android/iOS).

## Требования

- Flutter SDK 3.0+
- Запущенный бэкенд API (см. корень репозитория)

## Запуск

1. Из корня репозитория поднять БД и API:
   ```bash
   docker compose up -d
   go run ./cmd/api
   ```
2. Из папки `nc_app`:
   - **Веб:** `flutter run -d chrome`
   - **Android:** `flutter run -d android` (для эмулятора API: `flutter run -d android --dart-define=NC_API_URL=http://10.0.2.2:8080`)
   - **iOS:** `flutter run -d ios`

## Конфигурация API

По умолчанию используется `http://localhost:8080`. Для эмулятора Android укажите:
`--dart-define=NC_API_URL=http://10.0.2.2:8080`

## Тестовые аккаунты

Пароль у всех: **password123**

- student1@nc.kz (студент)
- teacher1@nc.kz (преподаватель)
- director@nc.kz (директор)
- admin@nc.kz (админ)
