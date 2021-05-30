package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/dumbogo/chess/engine"
	"github.com/dumbogo/chess/messagebroker"
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
	Email     string `gorm:"unique"`
	Name      string
	FirstName string
	LastName  string
	NickName  string
	UserID    string
	// TODO: hash AccessToken, AccessTokenSecret, RefreshToken
	AccessToken       string
	AccessTokenSecret string
	RefreshToken      string
	ExpiresAt         sql.NullTime
	IDToken           string
}

// GetUserFromAccessToken returns user from database with accesstoken set
func GetUserFromAccessToken(accessToken string) User {
	// TODO: Salt AccessToken
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
	Movements     []Movement // TODO: this will cause problems, implement when its needed by the engine
	WhitePieces   pieces     `gorm:"type:jsonb;not null"`
	BlackPieces   pieces     `gorm:"type:jsonb;not null"`
	BoardSquares  Squares    `gorm:"type:jsonb;not null"`
}

// AfterSave ...
func (g *Game) AfterSave(tx *gorm.DB) (err error) {
	if MessageBroker == nil {
		return nil
	}
	board := engine.LoadBoard(&engine.Player{}, &engine.Player{}, squaresToEngineSquares(g.BoardSquares))
	bytes, err := json.Marshal(payloadUpdateGame{
		Turn:   fmt.Sprint(g.Turn), // TODO: add corresponding color player
		Status: "somestatus",       // TODO: add the corresponding status
		Board:  board.String(),
	})
	if err != nil {
		return err
	}
	MessageBroker.Publish(g.UUID.String(), messagebroker.Message{Payload: bytes})
	return
}

type payloadUpdateGame struct {
	Turn   string `json:"turn"`
	Board  string `json:"board"`
	Status string `json:"status"`
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

// Squares Game squares stored in Model
type Squares map[engine.SquareIdentifier]Square

// Square square
type Square struct {
	engine.Square
	Piece Piece
}

func squaresToEngineSquares(bs Squares) engine.Squares {
	squares := engine.Squares{}
	for i, v := range bs {
		sq := engine.Square{
			Empty:            v.Empty,
			Coordinates:      v.Coordinates,
			SquareIdentifier: v.SquareIdentifier,
		}
		if !v.Empty {
			sq.Piece = engine.PieceFromPieceIdentifier(v.Piece.PieceIdentifier, v.Piece.Color)
		}
		squares[i] = sq
	}
	return squares
}

func engineSquaresToSquares(es engine.Squares) Squares {
	newSquares := Squares{}
	for i, v := range es {
		sq := newBasicSquare(v)
		if !sq.Empty {
			sq.Piece = Piece{
				PieceIdentifier: v.Piece.Identifier(),
				Color:           v.Piece.Color(),
			}
		}
		newSquares[i] = sq
	}
	return newSquares
}

// Piece piece
type Piece struct {
	PieceIdentifier engine.PieceIdentifier
	Color           engine.Color
}

// Scan implement scanner intarface
func (s *Squares) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	var result map[string]interface{}
	err := json.Unmarshal(bytes, &result) // returns not nil error
	if err != nil {
		return err
	}

	squares := make(Squares)

	for i, v := range result {
		idx, _ := strconv.Atoi(i)
		mapV := v.(map[string]interface{})
		squareid := mapV["SquareIdentifier"].(float64)
		mapCoordinates := mapV["Coordinates"].(map[string]interface{})
		xval, _ := mapCoordinates["X"].(float64)
		yval, _ := mapCoordinates["Y"].(float64)
		square := Square{
			engine.Square{
				Empty:            mapV["Empty"] == true,
				SquareIdentifier: engine.SquareIdentifier(squareid),
				Coordinates: engine.Coordinate{
					X: uint8(xval),
					Y: uint8(yval),
				},
			},
			Piece{},
		}
		if !square.Empty {
			pMap := mapV["Piece"].(map[string]interface{})
			pieceid := pMap["PieceIdentifier"].(float64)
			color := pMap["Color"].(float64)
			square.Piece = Piece{
				PieceIdentifier: engine.PieceIdentifier(pieceid),
				Color:           engine.Color(color),
			}
		}
		squares[engine.SquareIdentifier(idx)] = square
	}

	*s = squares
	return nil
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
