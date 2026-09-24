# gosdk-redis-core

`gosdk-redis-core` — пакет для удобной и безопасной работы с Redis через **go-redis/v9** в рамках  `gosdk-core` приложения.

- [Документация GOSDK-CORE](https://github.com/exgamer/gosdk-core)

- 🧩 **Dependency Injection**
  - [Что доступно в DI из коробки](pkg/di/DI_FUNCTIONS_README.MD)

  - [Использование REDIS HELPER](pkg/redis/REDIS_HELPER_README.MD)
## Возможности

- Корректное закрытие всех подключений
- Kernel для жизненного цикла приложения
- Helper-функции для бизнес-кода

## Установка

```bash
go get github.com/exgamer/gosdk-redis-core
```

## Переменные окружения

| Переменная | Описание |
|---|---|
| `REDIS_HOST` | адрес с портом, `localhost:6379` |
| `REDIS_USER` | пользователь ACL (Redis 6+). Не задан — вход только по паролю |
| `REDIS_PASSWORD` | пароль |
| `REDIS_DB` | номер базы |
| `REDIS_POOL_SIZE` | размер пула соединений. По умолчанию 10 на CPU |
| `REDIS_IDLE_TIMEOUT` | сколько секунд соединение может простаивать в пуле. По умолчанию 30 мин |
| `REDIS_MAX_CONN_LIFETIME` | максимальное время жизни соединения, секунды. По умолчанию без ограничения |
| `REDIS_DIAL_TIMEOUT_MS` | таймаут подключения, мс. По умолчанию 5000 |
| `REDIS_READ_TIMEOUT_MS` | таймаут чтения ответа, мс. По умолчанию 3000 |
| `REDIS_WRITE_TIMEOUT_MS` | таймаут записи, мс. По умолчанию равен таймауту чтения |
| `REDIS_MAX_RETRIES` | число повторов при сетевой ошибке. По умолчанию 3, `-1` — без повторов |

Не заданные значения (или `0`) — значения go-redis по умолчанию.

С умолчаниями зависший Redis держит каждый вызов до 12–15 с (3 с на чтение и повторы). Сервисам, которые должны быстро деградировать при недоступном Redis, стоит задать короткие значения, например:

```
REDIS_DIAL_TIMEOUT_MS=1000
REDIS_READ_TIMEOUT_MS=500
REDIS_WRITE_TIMEOUT_MS=500
REDIS_MAX_RETRIES=1
```

При старте конфиг выводится в stdout, пароль скрывается.

## Быстрый старт

### Подключение kernel

```go
a := app.NewApp()
_ = a.RegisterKernel(&app.RedisKernel{})
```

### Получение подключения

```go
db, err := app.GetRedisClient(a)
if err != nil {
    return err
}
```

## Shutdown

При остановке приложения автоматически вызывается `Close()` и все соединение закрывается.

## License

MIT
