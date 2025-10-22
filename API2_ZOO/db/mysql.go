package db

import (
	"api2/src/entities"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func Connect() {
	url := os.Getenv("MYSQLCONN")
	if url == "" {
		log.Fatal("MYSQLCONN environment variable is not set")
	}
	dsn := url
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatal("Database connection failed", err)
	}

	db.AutoMigrate(&entities.Visitas{}, &entities.Atraccion{}, &entities.VisitaGeneral{})
	DB = db
}
