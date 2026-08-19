REST API для управления контактами на Go. 
Поддерживает создание, чтение, обновление и удаление контактов, а также фильтрацию по избранным.
Технологии:
**Go** (net/http, gorilla/mux)
**PostgreSQL** (драйвер pgx)
**Docker** и **Docker Compose**
**Миграции** (golang-migrate)
**Логирование** (zap)
Запуск:
локально-
1. Скопируй `.env.example` в `.env` и заполни переменные
2. Запусти PostgreSQL
3. Выполни `make service-run`
В докер контейнерах-
Выполни `docker compose build`, затем `make docker-upd` (detached mod), `make docker-up` (не detached)
Примеры запросов: 
Приложение будет доступно по адресу http://localhost:9091
Создать контакт: POST /contacts Тело запроса json
{
  "name": "name",
  "num": "number",
  "isfav" true/false
}
Получить список всех контактов: GET /contacts
Получить избранные контакты: GET /contacts/fav
Изменить контакт: PUT /contacts/(id выбранного контакта) тело:
{
  "name": "name",
  "num": "number",
  "isfav" true/false
}
Удалить контакт: DELETE /contacts/(id выбранного контакта)
