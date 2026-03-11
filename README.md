# gosdk-redis-core

`gosdk-redis-core` — пакет для удобной и безопасной работы с Redis через **go-redis/v9** в рамках  `gosdk-core` приложения.

- [Документация GOSDK-CORE](https://github.com/exgamer/gosdk-core)

- 🧩 **Dependency Injection**
  - [Что доступно в DI из коробки](pkg/di/DI_FUNCTIONS_README.MD)

  - [Использование REDIS HELPER](pkg/redishelper/REDIS_HELPER_README.MD)
## Возможности

- Корректное закрытие всех подключений
- Kernel для жизненного цикла приложения
- Helper-функции для бизнес-кода

## Логирование запросов

Для вывода логов запросов используется ENV POSTGRES_DB_LOG_LEVEL, если не указан логи выключены
- "info"
- "errors"
- "warnings"


## Установка

```bash
go get github.com/exgamer/gosdk-redis-core
```

REDIS_HOST=localhost:6379
REDIS_DB=1

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
