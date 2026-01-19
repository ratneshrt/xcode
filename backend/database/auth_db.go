package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var AuthDB *gorm.DB

func ConnectAuthDB() {
	var err error
	dsn := "postgresql://postgres:mypassword@localhost:5432/postgres?sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	AuthDB = db
	log.Println("AuthDB connected")
}
