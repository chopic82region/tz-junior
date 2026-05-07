# tz-junior — тестовое задание (Go + PostgreSQL)

Небольшой REST API для управления пользователями и их подписками (создание/просмотр/обновление/удаление пользователей, создание/просмотр/отмена подписок, расчёт суммарной стоимости по фильтру).

## Что это за проект

Это решение **тестового задания по ТЗ на вакансию**. Проект реализует простой HTTP API и работу с PostgreSQL, включая миграции.

## Технологии

- **Go** (модуль: `github.com/chopic82region/tz-junior.git`)
- **HTTP framework**: `gin-gonic/gin`
- **PostgreSQL драйвер**: `jackc/pgx` (через `database/sql`)
- **Миграции**: `golang-migrate/migrate`
- **UUID**: `google/uuid`

## Структура проекта

- `cmd/main.go`: точка входа, загрузка конфигурации, подключение к БД, запуск миграций, старт HTTP-сервера
- `internal/config/config.go`: конфигурация (env/flags, дефолты)
- `internal/models/models.go`: модели `User`, `Subscription`
- `internal/service/postgres/postgres.go`: репозитории работы с PostgreSQL (users/subscriptions/filter)
- `internal/transport/handlers/handlers.go`: HTTP-хендлеры
- `internal/transport/server/server.go`: роутинг и запуск сервера
- `migrations/`: SQL миграции схемы
- `docs/swagger/`: Swagger/OpenAPI документация

## Конфигурация

Конфиг читается из **флагов** или **переменных окружения** (флаги имеют приоритет).

Обязательное:

- **`DB_PASSWORD`**: пароль PostgreSQL пользователя (без него приложение не стартует)

Остальное (есть дефолты):

- `DB_HOST` (default: `localhost`)
- `DB_PORT` (default: `5442`)
- `DB_USER` (default: `postgres`)
- `DB_NAME` (default: `finance_db`)
- `DB_SSLMODE` (default: `disable`)
- `SERVER_PORT` (default: `8081`)
- `MIGRATIONS_PATH` (default: `migrations`)

`.env` файл **опционален** (удобно для локальной разработки).

## Запуск

1) Подними PostgreSQL и создай БД (по умолчанию `finance_db`).

2) Установи переменные окружения (минимум `DB_PASSWORD`), например в PowerShell:

```powershell
$env:DB_PASSWORD="postgres"
$env:DB_HOST="localhost"
$env:DB_PORT="5442"
$env:DB_USER="postgres"
$env:DB_NAME="finance_db"
$env:DB_SSLMODE="disable"
$env:SERVER_PORT="8081"
$env:MIGRATIONS_PATH="migrations"
```

3) Запусти приложение:

```powershell
go run ./cmd
```

При старте приложение:

- подключается к БД
- автоматически выполняет **миграции** (`migrations/`)
- поднимает HTTP сервер на `SERVER_PORT`

## Тестирование (go test)

Запуск всех тестов:

```powershell
go test ./...
```

В проекте есть unit-тесты на:

- ключевую логику валидации и дефолтов подписок (repo)
- корректные HTTP ответы (handlers)

## Swagger / OpenAPI

Swagger спецификация лежит в `docs/swagger/openapi.yaml`.
Инструкции: `docs/swagger/README.md`.

## Параметры для Postman (быстрый старт)

Base URL (пример): `http://localhost:8081`

Общий заголовок:

- `Content-Type: application/json`

### 1) POST /users — создать пользователя

- **Method**: POST
- **URL**: `{{baseUrl}}/users`
- **Body (raw JSON)**:

```json
{
  "name": "Alice",
  "email": "alice@example.com"
}
```

### 2) GET /users — список пользователей

- **Method**: GET
- **URL**: `{{baseUrl}}/users`

### 3) GET /users/:id — получить пользователя

- **Method**: GET
- **URL**: `{{baseUrl}}/users/{{userId}}`

### 4) PATCH /users/:id — обновить пользователя

- **Method**: PATCH
- **URL**: `{{baseUrl}}/users/{{userId}}`
- **Body (raw JSON)**:

```json
{
  "name": "Alice Updated",
  "email": "alice.updated@example.com"
}
```

### 5) DELETE /users/:id — удалить пользователя

- **Method**: DELETE
- **URL**: `{{baseUrl}}/users/{{userId}}`

### 6) POST /subscriptions — создать подписку (end_date опционален)

- **Method**: POST
- **URL**: `{{baseUrl}}/subscriptions`
- **Body (raw JSON)**:

Вариант без `end_date` (сервер поставит `payment_time + 1 month`):

```json
{
  "user_id": "{{userId}}",
  "service_name": "netflix",
  "price": "10.00",
  "payment_time": "2026-05-07T12:00:00Z"
}
```

Вариант с `end_date`:

```json
{
  "user_id": "{{userId}}",
  "service_name": "netflix",
  "price": "10.00",
  "payment_time": "2026-05-07T12:00:00Z",
  "end_date": "2026-06-07T12:00:00Z"
}
```

Важно:

- `price` должен быть **числом в строке** (например `"199.99"`), иначе будет 400
- `payment_time` и `end_date` — **RFC3339**

### 7) GET /subscriptions — список подписок

- **Method**: GET
- **URL**: `{{baseUrl}}/subscriptions`

### 8) GET /subscriptions/:id — получить подписку

- **Method**: GET
- **URL**: `{{baseUrl}}/subscriptions/{{subscriptionId}}`

### 9) DELETE /subscriptions/:id — отменить подписку

- **Method**: DELETE
- **URL**: `{{baseUrl}}/subscriptions/{{subscriptionId}}`

Ответ вернёт:

- `message`
- `subscription` (включая `end_date`)

### 10) GET /filter/total_cost — сумма по фильтру

- **Method**: GET
- **URL**:

`{{baseUrl}}/filter/total_cost?user_id={{userId}}&created_at=2026-05-01T00:00:00Z&end=2026-06-01T00:00:00Z&name=netflix`

Query params:

- `user_id` (UUID)
- `created_at` (RFC3339) — начало периода
- `end` (RFC3339) — конец периода
- `name` — название сервиса (равно `service_name`)

