package settings

import (
	"encoding/json"
	"fmt"

	"github.com/iogo-framework/logs"
)

type Redis struct {
	Host string
	Port int
}

func (config Config) Redis() (Redis, error) {
	var redis Redis

	err := json.Unmarshal(config.Components["redis"], &redis)
	if err != nil {
		logs.Warning("%s: %s", err.Error(), "missing or wrong 'redis' configuration, ignoring")
	}

	if redis.Host == "" {
		redis.Host = "redis"
		logs.Warning("missing redis 'host' configuration, assuming default value: 'redis'")
	}

	if redis.Port == 0 {
		redis.Port = 6379
		logs.Warning("missing redis 'port' configuration, assuming default value: 6379")
	}

	return redis, nil
}

func (r Redis) String() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}
