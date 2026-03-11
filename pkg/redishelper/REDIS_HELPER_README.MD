### Пример использования REDIS HELPER

#### Для отключения кеша можно отправить заголовок Cache-Control:no-cache и указать SetRequestData(appInfo)

## Записать в кеш массив
```go

    err := database.NewRedisHelper[models.Location](ctx, repo.redisClient).
		SetArray("banner-service:locations", mdls, 1*time.Hour)
    
    if err != nil {
        return err
    }
```

## Получить из кеша массив
```go
	result, err := database.NewRedisHelper[models.Location](ctx, repo.redisClient).
		GetArray("banner-service:locations")

    if err != nil {
        return nil, err
    }

```

## Так же есть функции SetModel, GetModel - получить и записать структуру по ключу; 
## SetString, GetString - записать и получить строку по ключу; 
## MSetModels, MGetModels - записать и получить несколько структур по ключу