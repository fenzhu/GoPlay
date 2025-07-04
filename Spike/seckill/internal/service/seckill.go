package service

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync"

	"seckill/internal/model"
	"seckill/internal/repository"
)

var (
	ErrOutOfStock      = errors.New("out of stock")
	ErrProductNotFound = errors.New("product not found")
)

// seckillScript is a Lua script to atomically check and decrement stock in Redis.
// KEYS[1]: product stock key (e.g., "product:1:stock")
// ARGV[1]: quantity to decrement (always 1 for seckill)
// Returns:
//
//	>= 0: new stock value (success)
//	-1:  stock not found (or key doesn't exist)
//	-2:  insufficient stock
const seckillScript = `
local stock_key = KEYS[1]
local quantity = tonumber(ARGV[1])

local current_stock = tonumber(redis.call('GET', stock_key))

if current_stock == nil then
    return -1
end

if current_stock >= quantity then
    return redis.call('DECRBY', stock_key, quantity)
else
    return -2
end
`

var (
	seckillScriptSHA  string
	seckillScriptOnce sync.Once
)

// getSeckillScriptSHA 保证只加载一次脚本，并返回 SHA1
func getSeckillScriptSHA() (string, error) {
	var err error
	seckillScriptOnce.Do(func() {
		sha, e := repository.RDB.ScriptLoad(context.Background(), seckillScript).Result()
		if e != nil {
			err = e
			return
		}
		seckillScriptSHA = sha
	})
	return seckillScriptSHA, err
}

// runSeckillScript 封装 EvalSha+Eval+ScriptLoad 的兜底逻辑
func runSeckillScript(productKey string) (int64, error) {
	sha, err := getSeckillScriptSHA()
	if err != nil {
		return 0, err
	}
	result, err := repository.RDB.EvalSha(context.Background(), sha, []string{productKey}, 1).Result()
	if err != nil {
		if strings.HasPrefix(err.Error(), "NOSCRIPT") {
			// Eval 兜底
			result, err = repository.RDB.Eval(context.Background(), seckillScript, []string{productKey}, 1).Result()
			if err != nil {
				return 0, err
			}
			// 刷新 SHA
			sha, shaErr := repository.RDB.ScriptLoad(context.Background(), seckillScript).Result()
			if shaErr == nil {
				seckillScriptSHA = sha
			}
		} else {
			return 0, err
		}
	}
	return result.(int64), nil
}

func Seckill(productID, userID int64) error {
	productKey := "product:" + strconv.FormatInt(productID, 10) + ":stock"

	scriptResult, err := runSeckillScript(productKey)
	if err != nil {
		return err
	}

	switch scriptResult {
	case -1:
		return ErrProductNotFound
	case -2:
		return ErrOutOfStock
	default:
		// Stock successfully decremented, scriptResult is the new stock value
		// Continue with order creation and MySQL stock update
	}

	// Create order and update MySQL stock in separate goroutines
	go func() {
		// Create order in DB
		_, err := repository.DB.Exec("INSERT INTO orders (product_id, user_id) VALUES (?, ?)", productID, userID)
		if err != nil {
			// Handle error, e.g., log it or try to requeue the order creation
		}
	}()

	go func() {
		// Decrement stock in MySQL
		_, err := repository.DB.Exec("UPDATE products SET stock = stock - 1 WHERE id = ? ", productID)
		if err != nil {
			// Handle error, e.g., log it
		}
	}()

	return nil
}

func SyncProductStockToRedis() error {
	products, err := repository.GetAllDBProducts()
	if err != nil {
		return err
	}

	for _, p := range products {
		key := "product:" + strconv.FormatInt(p.ID, 10)

		cmd := repository.RDB.HSet(context.Background(), key, "id", p.ID, "name", p.Name)
		if cmd.Err() != nil {
			return cmd.Err()
		}

		stockKey := key + ":stock"

		scmd := repository.RDB.Set(context.Background(), stockKey, p.Stock, 0)
		if scmd.Err() != nil {
			return scmd.Err()
		}

		repository.RDB.SAdd(context.Background(), "products", p.ID)
	}
	return nil
}

func GetAllProducts() ([]model.Product, error) {
	return repository.GetAllProducts()
}

func GetProduct(id string) (map[string]string, error) {
	return repository.GetProduct(id)
}

func ResetSystem() error {
	// 1. Clear orders table in MySQL
	_, err := repository.DB.Exec("TRUNCATE TABLE orders")
	if err != nil {
		return err
	}

	// 2. Reset product stock in MySQL (assuming initial stock is 10 for Sample Product)
	// For a more robust solution, you might store initial stock in a config or separate table.
	_, err = repository.DB.Exec("UPDATE products SET stock = 10 WHERE id = 1") // Assuming product ID 1 is the sample product
	if err != nil {
		return err
	}

	// 3. Sync product stock to Redis
	return SyncProductStockToRedis()
}
