package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

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
	cmd := client.Ping(context.Background())
	if cmd.Err() != nil {
		panic(cmd.Err())
	} else {
		fmt.Println("redis connected" + cmd.Val())
	}
	DB.redis = client
}

func GetRedis() *redis.Client {
	return DB.redis
}
