package database

import (
	"api1/src/entities"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func Connect() {
	dsn := "root:root@tcp(127.0.0.1:3306)/mydb?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatal("Database connection failed", err)
	}

	db.AutoMigrate(&entities.Visitas{}, &entities.Atraccion{})
	DB = db
}
