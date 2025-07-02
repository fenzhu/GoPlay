package service

import (
	"context"
	"errors"
	"strconv"

	"seckill/internal/model"
	"seckill/internal/repository"
)

var (
	ErrOutOfStock      = errors.New("out of stock")
	ErrProductNotFound = errors.New("product not found")
)

func Seckill(productID, userID int64) error {
	// 1. Atomically decrement stock in Redis using DECR
	// DECR returns the value of key after the decrement.
	newStock, err := repository.RDB.Decr(context.Background(), "product:"+strconv.FormatInt(productID, 10)+":stock").Result()
	if err != nil {
		// If key doesn't exist, it might indicate product not found or initial sync issue
		// For simplicity, we'll treat it as product not found here.
		return ErrProductNotFound
	}

	if newStock < 0 {
		// If stock goes below zero, it means oversold. Increment back and return error.
		_, err := repository.RDB.Incr(context.Background(), "product:"+strconv.FormatInt(productID, 10)+":stock").Result()
		if err != nil {
			// Log this error, as it indicates a problem with Redis rollback
			// For now, we'll just return the original out of stock error
		}
		return ErrOutOfStock
	}

	// 2. Create order and update MySQL stock in separate goroutines
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
	products, err := repository.GetAllProducts()
	if err != nil {
		return err
	}

	for _, p := range products {
		key := "product:" + strconv.FormatInt(p.ID, 10) + ":stock"
		val := strconv.FormatInt(p.Stock, 10)
		_, err := repository.RDB.Set(context.Background(), key, val, 0).Result()
		if err != nil {
			return err
		}
	}
	return nil
}

func GetAllProducts() ([]model.Product, error) {
	return repository.GetAllProducts()
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
