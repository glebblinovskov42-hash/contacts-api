REST API для управления контактами на Go. 
Поддерживает создание, чтение, обновление и удаление контактов, а также фильтрацию по избранным.
Технологии:
1. **Go** (net/http, gorilla/mux)
2. **PostgreSQL** (драйвер pgx)
3. **Docker** и **Docker Compose**
4. **Миграции** (golang-migrate)
5. **Логирование** (zap)
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
