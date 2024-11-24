package main

import (
	"fmt"
	"log"
	"os"
	"time"
    "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	quit = make(chan bool)
    dbHandle *gorm.DB
)

var user string
var password string
var db string
var host string
var port string
var ssl string

func init() {
	user = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	db = os.Getenv("POSTGRES_DB")
	host = os.Getenv("POSTGRES_HOST")
	port = os.Getenv("POSTGRES_PORT")
	ssl = os.Getenv("POSTGRES_SSL")
	dbHandle = initDb()
    log.Println(dbHandle)
}

func initDb() *gorm.DB{
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, db, port, ssl)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if (err != nil){
		log.Println("Database init error:", err)
	}
	return db
}

func main() {
    fmt.Println("Hello main")
    func() {
		for {
			select {
			case <-quit:
				log.Println("Quitting")
				return
			default:
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
}