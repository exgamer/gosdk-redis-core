# Changelog

## v1.0.15

Минимальная версия Go поднята до 1.27.3. API и поведение пакета не изменились.

### Изменено
- `go.mod`: `go 1.26.0` → `go 1.27.3`. Причина — уязвимости стандартной библиотеки go1.26.0
  (`govulncheck`), исправлены в новых версиях Go. На go1.27.3 `govulncheck` уязвимостей стандартной библиотеки не находит.
- `gosdk-core` поднят до v1.0.29 (переход на Go 1.27.3, API не изменился).

### ⚠️ Изменения поведения
- Потребителю нужен Go ≥ 1.27.3: после `go get` этой версии в его `go.mod` станет `go 1.27.3`.
  С `GOTOOLCHAIN=auto` toolchain скачается сам; с `GOTOOLCHAIN=local` на Go < 1.27.3 или в образе
  без доступа к `proxy.golang.org` сборка упадёт — образ сборки `golang:1.27.3`.

## v1.0.14

### Изменено
- `gosdk-core` поднят до v1.0.28: `logger.Error` пишет на уровне Error.

## v1.1.0

### Добавлено
- Таймауты и повторы клиента: `REDIS_DIAL_TIMEOUT_MS`, `REDIS_READ_TIMEOUT_MS`, `REDIS_WRITE_TIMEOUT_MS`, `REDIS_MAX_RETRIES`. Не заданы — значения go-redis по умолчанию, поведение не меняется.
- `app.NewRedisOptions(cfg)` — сборка `redis.Options` из конфига.
- `REDIS_CONTEXT_TIMEOUT_ENABLED` — передаётся в `redis.Options.ContextTimeoutEnabled`: вызов прерывается по дедлайну или отмене контекста. Не задан — `false`, поведение не меняется.

### Исправлено
- `REDIS_USER`, `REDIS_POOL_SIZE`, `REDIS_IDLE_TIMEOUT`, `REDIS_MAX_CONN_LIFETIME` читались, но не передавались клиенту. Теперь применяются (`IDLE_TIMEOUT` и `MAX_CONN_LIFETIME` — в секундах). `REDIS_USER` передаётся только если задан явно.
- Пароль Redis больше не выводится в лог при старте.
