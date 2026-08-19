# Contacts API

REST API для управления контактами на Go.  
Поддерживает создание, чтение, обновление и удаление контактов, а также фильтрацию по избранным.

---

## Технологии

- **Go** (net/http, gorilla/mux)
- **PostgreSQL** (драйвер pgx)
- **Docker** и **Docker Compose**
- **Миграции** (golang-migrate)
- **Логирование** (zap)

---

### Локально

1. Скопируй `.env.example` в `.env` и заполни переменные
2. Запусти PostgreSQL
3. Выполни команду:

make service-run

docker compose build
make docker-upd   # detached mode
# или
make docker-up    # без detached

#### Примеры запросов
Создать контакт
POST /contacts
Content-Type: application/json

{
    "name": "Alice",
    "num": "12345678901",
    "isfav": true
}

Получить все контакты:
GET /contacts
Получить избранные контакты:
GET /contacts/fav
Изменить контакт:
PUT /contacts/{id}
Content-Type: application/json

{
    "name": "Bob",
    "num": "98765432100",
    "isfav": false
}
Удалить контакт:
DELETE /contacts/{id}
