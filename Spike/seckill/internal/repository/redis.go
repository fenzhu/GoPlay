package repository

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

func InitRedis(config map[string]interface{}) error {
	redisConfig := config["redis"].(map[string]interface{})
	addr := fmt.Sprintf("%s:%v", redisConfig["host"], redisConfig["port"])

	RDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: fmt.Sprintf("%v", redisConfig["password"]),
		DB:       0, // use default DB
		PoolSize: runtime.NumCPU() * 10, // Connection pool size, adjust based on your needs
		MinIdleConns: runtime.NumCPU() * 2, // Minimum number of idle connections
		PoolTimeout: 5 * time.Second, // Amount of time client waits for connection if all connections are busy
		ReadTimeout: 3 * time.Second, // Read timeout
		WriteTimeout: 3 * time.Second, // Write timeout
	})

	_, err := RDB.Ping(context.Background()).Result()
	return err
}