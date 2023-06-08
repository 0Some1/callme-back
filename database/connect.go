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
	dnsPostgres := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=verify-full",
		lib.DB_USER,
		lib.DB_PASSWORD,
		lib.DB_HOST,
		lib.DB_PORT,
		lib.DB_NAME,
	)

	database, err := gorm.Open(postgres.Open(dnsPostgres), &gorm.Config{})
	if err != nil {
		panic("could not connect to database")
	}
	DB.db = database
	migration(DB.db)
}

func migration(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Request{},
		&models.Comment{},
		&models.Photo{},
	)
	if err != nil {
		panic("could not migrate to database")
	}

}
