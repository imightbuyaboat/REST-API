# REST-API

Это простой RESTful API для управления сущностями Tasks, написанный на языке Go, с использованием модульной архитектуры и примерами взаимодействия с базой данных и кэшем.

## Возможности

- CRUD API на Go с использованием `net/http` и `github.com/gorilla/mux`
- Структурированные модули: обработчики, кэш, база данных, типы
- Использования Redis в качестве кэша
- Unit-тесты с использованием стандартного пакета `testing` и `mock`

## Описание модулей

1. `basic_types` - базовые типы для взаимодействия с базой данных и кэшем;
2. `db` - модуль для взаимодействия с базой данных `Postgress`;
3. `cache`- модуль для взаимодействия с `Redis`.

## Требования

- Go версии 1.16 и выше
- Docker и Docker-compose

## Установка и запуск

1. Клонируйте репозиторий:

   ```bash
   git clone https://github.com/imightbuyaboat/REST-API
   cd REST-API
   ```
   
2. В корне проекта создайте `.env` файл

   ```bash
   nano .env
   ```

   со следующим содержимым:

   ```env
   SQL_HOST=localhost
   SQL_PORT=5432
   SQL_DB=your_data_base
   SQL_USER=your_user
   SQL_PASSWORD=your_password

   REDIS_HOST=localhost
   REDIS_PORT=6379
   REDIS_PASSWORD=your_password
   ```
   
4. Запустите контейнеры через Docker-compose
   ```bash
   docker-compose up --build -d
   ```

5. Установите зависимости
   ```bash
   go mod download
   ```

6. Запустите веб-сервер
   ```bash
   go run .
   ```

   При успешном запуске в консоль будет выведено сообщение `Starting server at :8080`.

## Использование
API предоставляет базовые операции с сущностями Task (поля: id, name, description)

### POST (создание сущности Task c id = 1)

   ```bash
   curl -X POST http://localhost:8080/tasks/1 \
   -H "Content-Type: application/json" \
   -d '{"name": "Item 1", "description": "A test item"}'
   ```

### GET (получение сущности Task c id = 1)

   ```bash
   curl -X GET http://localhost:8080/tasks/1
   ```

### GET (получение всех сущностей Task)

   ```bash
   curl -X GET http://localhost:8080/tasks
   ```

### PUT (обновление всех полей сущности Task c id = 1)

   ```bash
   curl -X PUT http://localhost:8080/tasks/1 \
   -H "Content-Type: application/json" \
   -d '{"name": "Updated item", "description": "Updated description"}'
   ```

### DELETE (удаление сущности Task c id = 1)

   ```bash
   curl -X DELETE http://localhost:8080/tasks/1
   ```

## Тестирование

Для запуска всех тестов выполните:

   ```bash
   go test .
   ```
