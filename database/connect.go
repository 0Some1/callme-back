package database

import (
	"callme/models"
	"callme/utilities"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {

	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		utilities.Config("DB_USER"),
		utilities.Config("DB_PASSWORD"),
		utilities.Config("DB_HOST"),
		utilities.Config("DB_PORT"),
		utilities.Config("DB_NAME"))

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
