package models

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Entry struct {
	ID        uint
	Text      string
	Checked   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

var DB *gorm.DB

func ConnectDatabase() {
	db, err := gorm.Open(sqlite.Open("shopping_list.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	db.AutoMigrate(&Entry{})

	DB = db
}
