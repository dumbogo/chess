package api

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/dumbogo/chess/engine"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDbConn initialize database connection postgresql
func InitDbConn(host, port, user, password, dbname string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// Migrate runs migrations
func Migrate(db *gorm.DB) error {
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto;").Error; err != nil {
		log.Fatalf("Failed to create extension pgcrypto, got error %v", err)
	}
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		log.Fatalf("Failed to create extension uuid-ossp, got error %v", err)
	}
	return db.AutoMigrate(
		&User{},
		&Game{},
		&Player{},
		&Movement{},
	)
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
	Name          string
	UUID          uuid.UUID `gorm:"type:uuid;not null;default:uuid_generate_v1()"`
	WhitePlayer   Player
	WhitePlayerID sql.NullInt32
	BlackPlayer   Player
	BlackPlayerID sql.NullInt32
	Turn          int `gorm:"default:1"`
	Winner        int
	Movements     []Movement
	WhitePieces   map[uint8]uint8 `gorm:"type:jsonb;not null"`
	BlackPieces   map[uint8]uint8 `gorm:"type:jsonb;not null"`
	BoardSquares  engine.Squares  `gorm:"type:jsonb;not null"`
}

// Player Model
type Player struct {
	gorm.Model
	Color  string
	User   User
	UserID uint
}

// Movement Model
type Movement struct {
	gorm.Model
	PieceMoved int
	PieceEaten int
	From       string
	To         string
	Player     Player
	PlayerID   uint
	Game       Game
	GameID     uint
}
