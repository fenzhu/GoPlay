package database

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

type DataCenter struct {
	databases map[string]*sql.DB
	caches    map[string]*Cache
}

var Center *DataCenter = &DataCenter{
	databases: make(map[string]*sql.DB),
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

func (d *DataCenter) CreateDatabase(option *DataOption) (*sql.DB, error) {
	cfg := mysql.Config{
		User:                 option.User,
		Passwd:               option.Passwd,
		Net:                  "tcp",
		Addr:                 option.Addr,
		DBName:               option.DBName,
		AllowNativePasswords: true,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
		return nil, err
	}

	// fmt.Printf("%s connected\n", option.DBName)
	d.databases[option.DBName] = db
	return db, nil
}

func (d *DataCenter) GetDatabase(name string) *sql.DB {
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
