##  Что доступно из DI после инициализации приложения

Получить клиент redis:

```go
redisClient, err := di.GetRedisClient(c *di.Container) (*redis.Client, error)
```