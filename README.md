# URL Shortener
Сервис сокращения ссылок на Go.

## Подход к решению
Решение выполнено на стандартной библиотеке `net/http` 
без внешнего роутера, чтобы контролировать всю HTTP-цепочку (router, middleware, handlers) 
и явно обрабатывать ошибки.

## Архитектура
Проект построен по трехслойной архитектуре:
- `internal/storage` — реализации хранилища (`inmemory`, `postgres`).
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
- `400` — невалидный alias (длина/алфавит)
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

Скопировать переменные окружения:
```bash
cp .env.example .env
```

```bash
docker compose config
```
```bash
docker compose up --build
```

Что делает `docker compose up --build`:
- поднимает `db` (PostgreSQL);
- запускает `migrate` и применяет SQL-миграции через `goose`;
- после успешных миграций запускает `app`.

#### Переключение storage
По умолчанию контейнер запускается с `--storage=postgres` (в `Dockerfile` через `CMD`)

Можно изменить на `--storage=inmemory`

Сервис будет доступен на http://localhost:8080
## Тесты

```bash
go test ./...
```

```bash
go test -race ./...
```

Интеграционные тесты PostgreSQL:
```bash
docker compose up -d db
go test -tags=integration ./internal/storage/postgres -v
```

## Алгоритм создания alias
- генерируется массив длиной 10 байт с криптографически случайными значениями;
- принимаются только байты < 252 т.к. 252 кратно 63;
- байты 252..255 отбрасываются и генерируются заново;
- индекс символа вычисляется как byte % 63;
- итоговый alias состоит только из разрешенных символов;
- пространство значений: 63^10 комбинаций.

## Почему алгоритм устойчив
- в БД есть уникальность по `alias` (PK) и по `original_url` (UNIQUE)
- при `alias conflict` сервис генерирует новый alias и повторяет попытку;
- при повторном `original_url` PostgreSQL-слой атомарно возвращает уже существующий alias в одном SQL-запросе (`INSERT ... ON CONFLICT ... RETURNING` + `COALESCE`);
- это сохраняет инвариант: один `original_url` -> один alias.

## Правила URL
- принимаются только URL со схемой `http` или `https`;
- перед сохранением URL нормализуется:
  - удаляются внешние пробелы;
  - схема и host приводятся к нижнему регистру;
  - удаляются default ports (`:80` для `http`, `:443` для `https`);
  - если path пустой, устанавливается `/`.
