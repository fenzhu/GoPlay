package service

import (
	"context"
	"time"

	"example.com/redisweb/database"
	"github.com/redis/go-redis/v9"
)

/*
login HSet

	token: user

recent ZSet

	token: timestamp

viewed:token ZSet

	item: timestamp
*/
func client() *redis.Client {
	return database.GetRedis()
}

var (
	cli = client()
	ctx = context.Background()
	//set
	LOGIN = "login"
	//zset
	RECENT = "recent"
	//zset
	VIEWED_PREFIX = "viewed:"
)

func CheckToken(token string) string {
	cmd := cli.HGet(ctx, LOGIN, token)
	return cmd.Val()
}

func UpdateToken(token string, user string, item string) {
	client := client()

	timestamp := time.Now().Unix()
	client.HSet(ctx, LOGIN, token, user)
	client.ZAdd(ctx, RECENT, redis.Z{
		Member: token,
		Score:  float64(timestamp),
	})

	if len(item) > 0 {
		viewKey := VIEWED_PREFIX + token
		client.ZAdd(ctx, viewKey, redis.Z{
			Member: item,
			Score:  float64(timestamp),
		})
		//remain latest 25 items
		client.ZRemRangeByRank(ctx, viewKey, 0, -26)
	}
}
