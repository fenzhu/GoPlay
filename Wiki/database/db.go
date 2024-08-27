package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

type DataCenter struct {
	databases map[string]*sql.DB
}

var Center *DataCenter = &DataCenter{
	databases: make(map[string]*sql.DB),
}

type DataOption struct {
	User   string
	Passwd string
	Addr   string
	DBName string
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

	fmt.Printf("%s connected\n", option.DBName)
	d.databases[option.DBName] = db
	return db, nil
}

func (d *DataCenter) GetDatabase(name string) *sql.DB {
	return d.databases[name]
}
