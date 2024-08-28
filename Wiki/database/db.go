package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DataCenter struct {
	databases map[string]*gorm.DB
	caches    map[string]*Cache
}

var Center *DataCenter = &DataCenter{
	databases: make(map[string]*gorm.DB),
	caches:    make(map[string]*Cache),
}

type DataOption struct {
	User   string
	Passwd string
	Addr   string
	DBName string
}

type CacheOption struct {
	Name string
}

func (d *DataCenter) CreateDatabase(option *DataOption) (*gorm.DB, error) {
	// cfg := mysql.Config{
	// 	User:                 option.User,
	// 	Passwd:               option.Passwd,
	// 	Net:                  "tcp",
	// 	Addr:                 option.Addr,
	// 	DBName:               option.DBName,
	// 	AllowNativePasswords: true,
	// }
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		option.User, option.Passwd, option.Addr, option.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// gorm.Open(sql.Open("mysql"))
	// db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// db.
	// pingErr := db.Ping()
	// if pingErr != nil {
	// 	log.Fatal(pingErr)
	// 	return nil, err
	// }

	// fmt.Printf("%s connected\n", option.DBName)
	d.databases[option.DBName] = db
	return db, nil
}

func (d *DataCenter) GetDatabase(name string) *gorm.DB {
	return d.databases[name]
}

func (d *DataCenter) CreateCache(option *CacheOption) (*Cache, error) {
	cache := &Cache{Data: make(map[string]string)}

	d.caches[option.Name] = cache
	return cache, nil
}

func (d *DataCenter) GetCache(name string) *Cache {
	return d.caches[name]
}
