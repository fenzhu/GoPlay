package database

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	file, err := os.Create("gorm-log.txt")
	if err != nil {
		panic(err)
	}

	newLogger := logger.New(
		log.New(file, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn, // Log level
			IgnoreRecordNotFoundError: false,       //
			ParameterizedQueries:      false,       // include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		option.User, option.Passwd, option.Addr, option.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic(err)
	}

	d.databases[option.DBName] = db
	return db, nil
}

func (d *DataCenter) GetDatabase(name string) *gorm.DB {
	return d.databases[name]
}

func (d *DataCenter) CreateCache(option *CacheOption) (*Cache, error) {
	cache := &Cache{Data: &sync.Map{}}

	d.caches[option.Name] = cache
	return cache, nil
}

func (d *DataCenter) GetCache(name string) *Cache {
	return d.caches[name]
}
