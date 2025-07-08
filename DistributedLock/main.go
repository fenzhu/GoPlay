package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

var ctx = context.Background()

// RedisLock 定义了分布式锁的结构
type RedisLock struct {
	client *redis.Client
	key    string
	value  string
}

// NewRedisLock 创建一个新的分布式锁实例
func NewRedisLock(client *redis.Client, key string) *RedisLock {
	return &RedisLock{
		client: client,
		key:    key,
		value:  uuid.New().String(), // 为每个锁实例生成唯一值
	}
}

// Acquire 尝试获取锁
func (l *RedisLock) Acquire(expiration time.Duration) (bool, error) {
	ok, err := l.client.SetNX(ctx, l.key, l.value, expiration).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

// Release 释放锁
func (l *RedisLock) Release() error {
	// 使用 Lua 脚本确保原子性
	// 只有当 key 存在且 value 匹配时才删除
	script := `
        if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("del", KEYS[1])
        else
            return 0
        end
    `
	_, err := l.client.Eval(ctx, script, []string{l.key}, l.value).Result()
	return err
}

func main() {
	sync()
	// self()
}

func self() {
	// 初始化 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// 创建一个锁
	lockKey := "my-distributed-lock"
	lock := NewRedisLock(rdb, lockKey)

	// 尝试获取锁
	expiration := 10 * time.Second
	acquired, err := lock.Acquire(expiration)
	if err != nil {
		panic(err)
	}

	if acquired {
		fmt.Println("Lock acquired successfully!")

		// 模拟一些工作
		fmt.Println("Doing some work...")
		time.Sleep(5 * time.Second)

		// 释放锁
		err = lock.Release()
		if err != nil {
			panic(err)
		}
		fmt.Println("Lock released.")
	} else {
		fmt.Println("Could not acquire lock.")
	}
}
