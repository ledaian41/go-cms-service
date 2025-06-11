package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func Init(host string, user string, pwd string) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=5432 sslmode=disable", host, user, pwd)
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	//log.Println("Connected to:", version)
	//var err error
	//DB, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		log.Printf("❌ Failed to connect to database: %v", err)
	}
	log.Println("✅ Connected to database")

	sqlDB, _ := DB.DB()
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	return DB
}
