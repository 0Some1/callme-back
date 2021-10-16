package database

import (
	"callme/config"
	"callme/models"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {

	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.Config("DB_USER"),
		config.Config("DB_PASSWORD"),
		config.Config("DB_HOST"),
		config.Config("DB_PORT"),
		config.Config("DB_NAME"))

	database, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		panic("could not connect to database")
	}
	DB = database
	migration(DB)
}

func migration(db *gorm.DB) {
	db.AutoMigrate(&models.User{})
}
