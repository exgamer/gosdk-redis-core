# Changelog

## v1.1.0

### Добавлено
- Таймауты и повторы клиента: `REDIS_DIAL_TIMEOUT_MS`, `REDIS_READ_TIMEOUT_MS`, `REDIS_WRITE_TIMEOUT_MS`, `REDIS_MAX_RETRIES`. Не заданы — значения go-redis по умолчанию, поведение не меняется.
- `app.NewRedisOptions(cfg)` — сборка `redis.Options` из конфига.

### Исправлено
- `REDIS_USER`, `REDIS_POOL_SIZE`, `REDIS_IDLE_TIMEOUT`, `REDIS_MAX_CONN_LIFETIME` читались, но не передавались клиенту. Теперь применяются (`IDLE_TIMEOUT` и `MAX_CONN_LIFETIME` — в секундах). `REDIS_USER` передаётся только если задан явно.
- Пароль Redis больше не выводится в лог при старте.
