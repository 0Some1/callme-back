package database

import (
	"callme/lib"
	"callme/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB postgresDB

func Connect() {

	//dnsMysql := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
	//	lib.DB_USER,
	//	lib.DB_PASSWORD,
	//	lib.DB_HOST,
	//	lib.DB_PORT,
	//	lib.DB_NAME)

	dnsPostgres := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		lib.DB_HOST,
		lib.DB_USER,
		lib.DB_PASSWORD,
		lib.DB_NAME,
		lib.DB_PORT)

	database, err := gorm.Open(postgres.Open(dnsPostgres), &gorm.Config{})
	if err != nil {
		panic("could not connect to database")
	}
	DB.db = database
	migration(DB.db)
}

func migration(db *gorm.DB) {
	db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Request{},
		&models.Comment{},
		&models.Photo{},
	)

}
