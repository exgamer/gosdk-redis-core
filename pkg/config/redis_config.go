package config

// RedisConfig — настройки подключения к Redis из env.
//
// Таймауты и повторы: 0 (не задано) — значения go-redis по умолчанию
// (dial 5s, read/write 3s, 3 повтора). REDIS_MAX_RETRIES=-1 отключает повторы.
type RedisConfig struct {
	RedisHost     string `mapstructure:"REDIS_HOST" json:"redis_host"`
	RedisUser     string `mapstructure:"REDIS_USER" json:"redis_user"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD" json:"redis_password"`
	RedisDb       int    `mapstructure:"REDIS_DB" json:"redis_db"`

	PoolSize int `mapstructure:"REDIS_POOL_SIZE" json:"redis_pool_size"`
	// IdleTimeout — сколько секунд соединение может простаивать в пуле.
	IdleTimeout int `mapstructure:"REDIS_IDLE_TIMEOUT" json:"redis_idle_timeout"`
	// MaxConnLifetime — максимальное время жизни соединения, секунды.
	MaxConnLifetime int `mapstructure:"REDIS_MAX_CONN_LIFETIME" json:"redis_max_conn_lifetime"`

	DialTimeoutMs  int `mapstructure:"REDIS_DIAL_TIMEOUT_MS" json:"redis_dial_timeout_ms"`
	ReadTimeoutMs  int `mapstructure:"REDIS_READ_TIMEOUT_MS" json:"redis_read_timeout_ms"`
	WriteTimeoutMs int `mapstructure:"REDIS_WRITE_TIMEOUT_MS" json:"redis_write_timeout_ms"`
	MaxRetries     int `mapstructure:"REDIS_MAX_RETRIES" json:"redis_max_retries"`
}

// Masked возвращает копию конфига со скрытым паролем — для вывода в лог.
func (c RedisConfig) Masked() RedisConfig {
	if c.RedisPassword != "" {
		c.RedisPassword = "***"
	}

	return c
}
