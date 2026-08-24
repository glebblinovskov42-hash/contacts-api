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

### Для запуска локально

1. Скопируй `.env.example` в `.env` и заполни переменные
2. Запусти PostgreSQL
3. Выполни команду:

make service-run

---

#### Для применения миграций:
make migrate-up 
make migrate-down

---

#### Для запуска в докер-контейнерах

docker compose build
make docker-upd   # detached mode
или
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

---

Получить избранные контакты:
GET /contacts/fav
(Также для получения всех контактов и избранных контактов доступна пагинация через query параметры: ?page={int}&limit={int})

---

Изменить контакт:
PUT /contacts/{id}
Content-Type: application/json

{
    "name": "Bob",
    "num": "98765432100",
    "isfav": false
}

---

Удалить контакт:
DELETE /contacts/{id}

---

##### Для запуска unit-тестов:

make unit-test

###### Для запуска e2e-теста:

make e2e-test
