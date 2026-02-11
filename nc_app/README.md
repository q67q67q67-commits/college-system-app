# NC App — мобильное приложение Narxoz College

Flutter-приложение для студентов и сотрудников колледжа. Подключается к Go API проекта NC.

## Требования

- [Flutter SDK](https://flutter.dev/docs/get-started/install) (3.0+)
- Запущенный бэкенд: `go run ./cmd/api` из корня репозитория (и PostgreSQL через `docker compose up -d`)

## Первый запуск

1. Из корня репозитория перейдите в папку приложения:
   ```bash
   cd nc_app
   ```
2. Создайте платформенные проекты (Android/iOS/Web), если папок `android/`, `ios/` ещё нет:
   ```bash
   flutter create .
   ```
3. Установите зависимости:
   ```bash
   flutter pub get
   ```
4. Запустите приложение:
   - Эмулятор или устройство: `flutter run`
   - Chrome: `flutter run -d chrome`

## Настройка URL API

По умолчанию используется `http://localhost:8080`. Для Android-эмулятора замените в коде на `http://10.0.2.2:8080`, для реального устройства — на IP вашего компьютера в сети (например `http://192.168.1.100:8080`). Константа находится в `lib/src/services/api_client.dart`: `defaultBaseUrl`.

## Тестовый вход

Используйте учётные данные из сидов бэкенда, например:
- Email: `student1@nc.kz`
- Пароль: `password123`

См. также `api/ENDPOINTS.md` и `context.md` в корне проекта.
