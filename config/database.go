package config

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/jake-schuler/fll-event-screen/models"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("data/event.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&models.Teams{})
	return db
}