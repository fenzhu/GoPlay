
package repository

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

func InitRedis(config map[string]interface{}) error {
	redisConfig := config["redis"].(map[string]interface{})
	addr := fmt.Sprintf("%s:%d", redisConfig["host"], redisConfig["port"])

	RDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: fmt.Sprintf("%v", redisConfig["password"]),
		DB:       0, // use default DB
	})

	_, err := RDB.Ping(context.Background()).Result()
	return err
}
