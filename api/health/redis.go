package health

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type RedisComponent struct {
	name string
	r    *redis.Client
}

func NewRedisComponent(name string, r *redis.Client) *RedisComponent {
	return &RedisComponent{
		name: name,
		r:    r,
	}
}

func (r *RedisComponent) HealthComponentName() string {
	return r.name
}

func (r *RedisComponent) CheckHealthy() (bool, string) {
	_, err := r.r.Ping(context.Background()).Result()
	if err != nil {
		return false, err.Error()
	}
	return true, "ok"
}
