package database

import "github.com/redis/go-redis/v9"

type Cache struct {
	Data *redis.Client
}
