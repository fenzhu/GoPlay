package service

import (
	"context"
	"time"

	"example.com/redisweb/database"
	"github.com/redis/go-redis/v9"
)

func client() *redis.Client {
	return database.GetRedis()
}

var (
	cli = client()
	ctx = context.Background()
)

func CheckToken(token string) string {

	cmd := cli.HGet(ctx, "login", token)
	return cmd.Val()
}

func UpdateToken(token string, user string, item *string) {
	client := client()

	timestamp := time.Now().Unix()
	client.HSet(ctx, "login", token, user)
	client.ZAdd(ctx, "recent", redis.Z{
		Member: token,
		Score:  float64(timestamp),
	})

	if item != nil {
		viewKey := "viewed:" + token
		client.ZAdd(ctx, viewKey, redis.Z{
			Member: *item,
			Score:  float64(timestamp),
		})
		//remain latest 25 items
		client.ZRemRangeByRank(ctx, viewKey, 0, -26)
	}
}
