package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"seckill/internal/model"
)

var DB *sql.DB

func InitDB(config map[string]interface{}) error {
	mysqlConfig := config["mysql"].(map[string]interface{})
	dsn := fmt.Sprintf("%s:%v@tcp(%s:%v)/%s",
		mysqlConfig["user"],
		fmt.Sprintf("%v", mysqlConfig["password"]),
		mysqlConfig["host"],
		mysqlConfig["port"],
		mysqlConfig["dbname"],
	)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	return DB.Ping()
}

func GetAllProducts() ([]model.Product, error) {
	rows, err := DB.Query("SELECT id, name, stock FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Stock); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}