package session

import (
	"context"
	"time"

	"example.com/redisweb/database"
	"example.com/redisweb/service"
	"github.com/redis/go-redis/v9"
)

var (
	cli         = client()
	ctx         = context.Background()
	quit        = false
	LIMIT int64 = 10000000
)

func client() *redis.Client {
	return database.GetRedis()
}

func CleanSessions() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for !quit {
		select {
		case <-ticker.C:
			cmd := cli.ZCard(ctx, service.RECENT)
			size := cmd.Val()
			if size <= LIMIT {
				continue
			}

			end_index := min(size-LIMIT, 100)
			zcmd := cli.ZRange(ctx, service.RECENT, 0, end_index-1)
			tokens := zcmd.Val()
			for _, token := range tokens {
				cli.HDel(ctx, service.LOGIN, token)
				cli.ZRem(ctx, service.RECENT, token)
				viewKey := service.VIEWED_PREFIX + token
				cli.Del(ctx, viewKey)
			}
		case <-ctx.Done():
			return
		}
	}
}
