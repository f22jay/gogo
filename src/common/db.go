package common

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

// USERNAME user
const (
	USERNAME = "root"
	PASSWORD = "passw0rd"
	NETWORK  = "tcp"
	SERVER   = "172.18.0.22"
	PORT     = 3306
	DATABASE = "my_db"
)

// RedisAddr : redis addr
const RedisAddr = "172.18.0.23:6379"

var db *sql.DB

func init() {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Open mysql failed,err:%v\n", err)
		return
	}
	db.SetConnMaxLifetime(100 * time.Second)
	db.SetMaxOpenConns(1000)
	db.SetMaxIdleConns(1000)
}

// GetDb : return db
func GetDb() *sql.DB {
	return db
}
