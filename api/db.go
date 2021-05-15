package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/dumbogo/chess/engine"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBConn stores DB Connection
var DBConn *gorm.DB

// InitDbConn initialize database connection to postgresql
func InitDbConn(host, port, user, password, dbname string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	var err error
	DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return DBConn, err
}

// Migrate runs migrations
func Migrate() error {
	if err := DBConn.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto;").Error; err != nil {
		log.Fatalf("Failed to create extension pgcrypto, got error %v", err)
	}
	if err := DBConn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		log.Fatalf("Failed to create extension uuid-ossp, got error %v", err)
	}
	return DBConn.AutoMigrate(
		&User{},
		&Game{},
		&Player{},
		&Movement{},
	)
}

// User Model
type User struct {
	gorm.Model
	Email             string `gorm:"unique"`
	Name              string
	FirstName         string
	LastName          string
	NickName          string
	UserID            string
	AccessToken       string
	AccessTokenSecret string
	RefreshToken      string
	ExpiresAt         sql.NullTime
	IDToken           string
}

// GetUserFromAccessToken returns user from database with accesstoken set
func GetUserFromAccessToken(accessToken string) User {
	user := User{}
	DBConn.Where("access_token=?", accessToken).First(&user)
	return user
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
	Turn          uint
	Winner        int
	Movements     []Movement // TODO: this will cause problems
	WhitePieces   pieces     `gorm:"type:jsonb;not null"`
	BlackPieces   pieces     `gorm:"type:jsonb;not null"`
	BoardSquares  squares    `gorm:"type:jsonb;not null"` // TODO: needs to change, not working, more info at line 79
}

type pieces map[uint8]uint8

// Scan scan value into Jsonb, implements sql.Scanner interface
func (p *pieces) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	result := pieces{}
	err := json.Unmarshal(bytes, &result)
	*p = result
	return err
}

// =-----------------EXPERIMENTAL=======================
type squares engine.Squares

// type square struct {
// 	engine.Square
// 	Piece piece `json:"piece"`
// }
//
// type piece struct {
// 	PieceIdentifier engine.PieceIdentifier `json:"piece_identifier"`
// 	Color           Color                  `json:"color"`
// }

// Scan implements Scanner interface
// TODO: Not decoding properly square.Piece property, returns error
func (s *squares) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	result := squares{}
	err := json.Unmarshal(bytes, &result) // returns not nil error
	// json.Unmarshal(bytes, &result)
	fmt.Printf("error: %+v\n", err)
	// error: json: cannot unmarshal object into Go struct field Square.Piece of type engine.Piece
	// *s = result
	return nil
}

// Value return json value, implement driver.Valuer interface
// func (s squares) Value() (driver.Value, error) {} // TODO: Implement

// =-----------------END EXPERIMENTAL=======================

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
