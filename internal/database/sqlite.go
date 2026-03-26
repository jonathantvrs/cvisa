package database

import (
	"log"

	"github.com/glebarez/sqlite"

	"github.com/jonathantvrs/cvisa/internal/model"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("credit_system.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	db.AutoMigrate(&model.AccountLimit{})
	return db
}
