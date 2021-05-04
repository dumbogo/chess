package api

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDb initialize database connection postgresql
func InitDb(host, port, user, password, dbname string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// Migrate runs migrations
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Game{}, &Player{})
}

// User Model
type User struct {
	gorm.Model
	Username  string
	PublicKey []byte
}

// Game Model
type Game struct {
	gorm.Model
}

// Player Model
type Player struct {
	gorm.Model
	Color string
}
