package database

import "github.com/redis/go-redis/v9"

type Option struct {
	Name string
	Addr string
}

type Database struct {
	redis *redis.Client
}

var DB *Database = &Database{}

func CreateRedis(option *Option) {
	client := redis.NewClient(&redis.Options{
		Addr:     option.Addr,
		Password: "",
		DB:       0,
	})

	DB.redis = client
}

func GetRedis() *redis.Client {
	return DB.redis
}
