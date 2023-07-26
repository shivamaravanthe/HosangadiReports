package database

import (
	"log"
	"shivamaravanthe/HosangadiReports/constants"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	DB, err = gorm.Open(mysql.Open(constants.DSN), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database %v", err)
		return
	}
	if DB.Error != nil {
		log.Printf("Failed to connect to database %v", err)
		return
	}
}
