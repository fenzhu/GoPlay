package repository

import (
	"context"
	"fmt"
	"runtime"
	"seckill/internal/model"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

func InitRedis(config map[string]interface{}) error {
	redisConfig := config["redis"].(map[string]interface{})
	addr := fmt.Sprintf("%s:%v", redisConfig["host"], redisConfig["port"])

	RDB = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     fmt.Sprintf("%v", redisConfig["password"]),
		DB:           0,                     // use default DB
		PoolSize:     runtime.NumCPU() * 10, // Connection pool size, adjust based on your needs
		MinIdleConns: runtime.NumCPU() * 2,  // Minimum number of idle connections
		PoolTimeout:  5 * time.Second,       // Amount of time client waits for connection if all connections are busy
		ReadTimeout:  3 * time.Second,       // Read timeout
		WriteTimeout: 3 * time.Second,       // Write timeout
	})

	_, err := RDB.Ping(context.Background()).Result()
	return err
}

func GetAllProducts() ([]model.Product, error) {
	cmd := RDB.SMembers(context.Background(), "products")

	res, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	var products []model.Product
	for _, id := range res {
		key := "product:" + id

		hCmd := RDB.HGetAll(context.Background(), key)
		if hCmd.Err() != nil {
			return nil, hCmd.Err()
		}

		productId, parseErr := strconv.ParseInt(id, 10, 64)
		if parseErr != nil {
			return nil, parseErr
		}

		val := hCmd.Val()
		p := model.Product{
			ID:   productId,
			Name: val["name"],
		}

		products = append(products, p)
	}

	return products, nil
}

func GetProduct(id string) (map[string]string, error) {
	key := "product:" + id

	hCmd := RDB.HGet(context.Background(), key, "name")
	if hCmd.Err() != nil {
		return nil, hCmd.Err()
	}

	val := hCmd.Val()

	m := map[string]string{"id": id, "name": val}

	return m, nil
}
