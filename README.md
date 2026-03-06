# URL Shortener
Сервис сокращения ссылок на Go

## Подход к решению
Решение выполнено на стандартной библиотеке `net/http` 
без внешнего роутера, чтобы контролировать всю HTTP-цепочку (router, middleware, handlers) 
и явно обрабатывать ошибки.

## Архитектура
Проект построен по трехслойной архитектуре:
- `internal/storage` — реализации хранилища (`in_memory`, `postgres`).
- `internal/services` — бизнес-логика (валидация URL, генерация alias, обработка конфликтов).
- `internal/http/handlers` — HTTP transport и маппинг ошибок в статус-коды.
- `internal/http/middleware` — method check, recover, logging.
- `internal/http/router` — сборка маршрутов и middleware chain.

## Требования к alias
- длина: `10` символов;
- алфавит: `[a-zA-Z0-9_]`;
- один `original_url` всегда маппится в один и тот же alias.

## API

### POST `/url`
Создает сокращенную ссылку или возвращает уже существующую для того же URL.

Request:
```json
{
  "url": "https://google.com"
}
```

Success (`201`):
```json
{
  "alias": "AbC123_xYz"
}
```

Error (пример):
```json
{
  "error": "invalid url"
}
```

Ошибки:
- `400` — невалидный JSON / невалидный URL
- `413` — слишком большой request body
- `500` — внутренняя ошибка
- `503` — превышено число попыток генерации alias

### GET `/url/{alias}`
Возвращает оригинальный URL по alias.

Success (`200`):
```json
{
  "url": "https://google.com"
}
```

Error (пример):
```json
{
  "error": "not found"
}
```

Ошибки:
- `400` — пустой alias
- `404` — alias не найден
- `500` — внутренняя ошибка

## Локальный запуск

```bash
go run ./cmd/app --storage=inmemory --config=./config/local.yaml
```

```bash
go run ./cmd/app --storage=postgres --config=./config/local.yaml
```

## Запуск через Docker

### 1. Поднять сервисы
```bash
docker compose up --build
```

### 2. Применить миграцию
```bash
docker compose exec -T db psql -U postgres -d postgres < migrations/migrations/000001_init.up.sql
```

### 3. Переключение storage
По умолчанию контейнер запускается с `--storage=postgres` (в `Dockerfile` через `CMD`)

Можно изменить на `--storage=inmemory`

Сервис будет доступен на http://localhost:8080
## Тесты

```bash
go test ./...
```

```bash
go test -race ./internal/storage/in_memory -run Concurrent
```

## Алгоритм создания alias
- генерируется массив длиной `10` байт с криптографически случайными значениями
- каждый байт маппится в символ алфавита через `byte % len(alphabet)`
- итоговый alias состоит только из разрешенных символов
- пространство значений: `63^10` комбинаций

## Почему алгоритм устойчив
- в БД есть уникальность по `alias` (PK) и по `original_url` (UNIQUE)
- при `alias conflict` сервис генерирует новый alias и повторяет попытку
- при `url conflict` (гонка между запросами) сервис читает уже созданный alias и возвращает его;
- это сохраняет инвариант: один `original_url` -> один alias.

