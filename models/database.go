package models

import (
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type List struct {
	Name      string `gorm:"primaryKey"`
	Password  []byte
	CreatedAt time.Time
	UpdatedAt time.Time
	Entries   []Entry
}

type Entry struct {
	ID        uint `gorm:"primaryKey"`
	Text      string
	Checked   bool
	CreatedAt time.Time
	UpdatedAt time.Time
	ListName  string
}

var DB *gorm.DB

func ConnectDatabase() {
	db, err := gorm.Open(sqlite.Open("shopping_list.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	db.AutoMigrate(&List{})
	db.AutoMigrate(&Entry{})

	DB = db
}

func BasePath() string {
	basePath := os.Getenv("BASE_PATH")
	// if basePath == "" {
		// basePath = "/"
	// }

	return basePath
}
