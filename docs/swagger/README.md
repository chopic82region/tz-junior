# Swagger / OpenAPI

Документация лежит в файле `openapi.yaml` (OpenAPI 3.0).

## Как открыть

- **Swagger Editor**: откройте файл `docs/swagger/openapi.yaml` в [Swagger Editor](https://editor.swagger.io/) (через *File → Import File*).
- **Postman / Insomnia**: импортируйте `openapi.yaml` как OpenAPI спецификацию.

## Как быстро проверить ручки

В `README.md` в корне проекта есть готовые параметры для Postman (URL/headers/body/query), раздел **“Параметры для Postman (быстрый старт)”**.

## Примечания

- Время в запросах/ответах — **RFC3339** (пример: `2026-05-07T12:00:00Z`).
- В `POST /subscriptions` поле `end_date` **опционально**: если не передано, сервер выставляет `payment_time + 1 month`.

