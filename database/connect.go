package database

import (
	"callme/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
func Connect()  {
	database , err:= gorm.Open(mysql.Open("root:4935@/test0"),&gorm.Config{})
	if err!=nil {
		panic("could not connect to database")
	}
	DB=database
	migration(DB)
}

func migration(db *gorm.DB )  {
	db.AutoMigrate(&models.User{})
}
