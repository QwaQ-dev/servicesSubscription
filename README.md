# Service Subscriptions — REST API
Краткое и наглядное описание REST-сервиса для агрегации данных об онлайн‑подписках.

## Что реализовано
- CRUDL для записей о подписках (Create, Read, Update, Delete, List).
- Эндпоинт для подсчёта суммарной стоимости подписок за период с фильтрацией по user_id и service_name.
- СУБД: PostgreSQL, миграции для инициализации БД.
- Логирование в stdout (конфигурируемое через .env).
- Конфигурация через `.env` / `.yaml`.
- Swagger‑документация.
- Запуск через Docker Compose.

---

## Работа с API

Базовый префикс: `/api/v1`

Эндпоинты:
- POST `/api/v1/subscriptions` — создать подписку
- GET `/api/v1/subscriptions/{id}` — получить подписку
- PUT `/api/v1/subscriptions/{id}` — обновить подписку
- DELETE `/api/v1/subscriptions/{id}` — удалить подписку
- GET `/api/v1/subscriptions` — список подписок
- GET `/api/v1/subscriptions/summ` — суммарная стоимость за период

Формат даты начала/окончания: `MM-YYYY` (пример: `07-2025`). Стоимость — целое число (рубли).

Пример тела запроса на создание:
````json
{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025"
}
