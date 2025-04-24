# REST-API

Это простой RESTful API, написанный на языке Go, с использованием модульной архитектуры и примерами взаимодействия с базой данных и кэшем.

## Возможности

- CRUD API на Go с использованием `net/http` и `github.com/gorilla/mux`
- Структурированные модули: обработчики, кэш, база данных, типы
- Использования Redis в качестве кэша
- Unit-тесты с использованием стандартного пакета `testing`

## Требования

- Go версии 1.16 и выше

## Установка и запуск

1. Клонируйте репозиторий:

   ```bash
   git clone https://github.com/imightbuyaboat/REST-API
   cd REST-API
   ```
   
2. Запустите веб-сервер:

   ```bash
   go run .
   ```

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
