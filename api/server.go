package api

import (
	context "context"
	"database/sql"
	"errors"

	"github.com/dumbogo/chess/engine"
	"gorm.io/gorm"
)

// Server grpc server interface implementation
type Server struct {
	UnimplementedChessServiceServer
	Db *gorm.DB
}

// StartGame starts a new game
func (s *Server) StartGame(ctx context.Context, startGameRequest *StartGameRequest) (*StartGameResponse, error) {
	// TODO: WIP, needs to get User from auth
	// 1.---------- change these lines to load user instead
	user := User{}
	s.Db.Create(&user)
	// 1.----------

	game := Game{
		Name: startGameRequest.GetName(),
		BlackPieces: map[uint8]uint8{
			uint8(engine.RookIdentifier):   2,
			uint8(engine.KnightIdentifier): 2,
			uint8(engine.BishopIdentifier): 2,
			uint8(engine.QueenIdentifier):  1,
			uint8(engine.KingIdentifier):   1,
			uint8(engine.PawnIdentifier):   8,
		},
		WhitePieces: map[uint8]uint8{
			uint8(engine.RookIdentifier):   2,
			uint8(engine.KnightIdentifier): 2,
			uint8(engine.BishopIdentifier): 2,
			uint8(engine.QueenIdentifier):  1,
			uint8(engine.KingIdentifier):   1,
			uint8(engine.PawnIdentifier):   8,
		},
		BoardSquares: squares(engine.PristineSquares()),
	}
	switch startGameRequest.GetColor() {
	case Color_WHITE:
		p := Player{
			Color:  Color.Enum(Color_WHITE).String(),
			UserID: user.ID,
		}
		s.Db.Create(&p)
		game.WhitePlayerID = sql.NullInt32{Valid: true, Int32: int32(p.ID)}
		game.Turn = p.ID
	case Color_BLACK:
		p := Player{
			Color:  Color.Enum(Color_BLACK).String(),
			UserID: user.ID,
		}
		s.Db.Create(&p)
		game.BlackPlayerID = sql.NullInt32{Valid: true, Int32: int32(p.ID)}
		game.Turn = p.ID
	}
	result := s.Db.Create(&game)
	if result.Error != nil {
		return nil, result.Error
	}

	startGameResponse := &StartGameResponse{
		Uuid: game.UUID.String(),
	}
	return startGameResponse, nil
}

// JoinGame joins a game
func (s *Server) JoinGame(ctx context.Context, r *JoinGameRequest) (*JoinGameResponse, error) {
	// TODO: WIP, needs to get User from auth
	// 1.---------- change these lines to load user instead
	user := User{}
	s.Db.Create(&user)
	// 1.----------

	uuid := r.GetUuid()

	var game Game
	tx := s.Db.Where("UUID = ?", uuid).First(&game)
	if tx.Error != nil {
		// TODO: Create a custom error instead of gorm error and return it
		return nil, tx.Error
	}

	if game.BlackPlayerID.Valid && game.WhitePlayerID.Valid {
		return nil, errors.New("already full game")
	}
	var color Color
	if game.BlackPlayerID.Valid {
		color = Color_WHITE
		player := Player{
			Color:  Color.Enum(Color_WHITE).String(),
			UserID: user.ID,
		}
		s.Db.Create(&player)
		if tx.Error != nil {
			return nil, tx.Error
		}
		game.WhitePlayerID = sql.NullInt32{Valid: true, Int32: int32(player.ID)}
	} else if game.WhitePlayerID.Valid {
		color = Color_BLACK
		player := Player{
			Color:  Color.Enum(Color_BLACK).String(),
			UserID: user.ID,
		}
		tx := s.Db.Create(&player)
		if tx.Error != nil {
			return nil, tx.Error
		}
		game.BlackPlayerID = sql.NullInt32{Valid: true, Int32: int32(player.ID)}
	}

	tx = s.Db.Save(&game)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &JoinGameResponse{
		Uuid:  uuid,
		Color: color,
	}, nil
}

// Move Moves a player piece
func (s *Server) Move(ctx context.Context, r *MoveRequest) (*MoveResponse, error) { // TODO: Implementation
	panic("Pending implementation")
}
