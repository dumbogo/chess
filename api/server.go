package api

import (
	context "context"
	"database/sql"

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
	user := User{}
	s.Db.Create(&user)

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
		BoardSquares: engine.PristineSquares(),
	}
	switch startGameRequest.GetColor() {
	case Color_WHITE:
		p := Player{
			Color:  Color.Enum(Color_WHITE).String(),
			UserID: user.ID,
		}
		s.Db.Create(&p)
		game.WhitePlayerID = sql.NullInt32{Valid: true, Int32: int32(p.ID)}
	case Color_BLACK:
		p := Player{
			Color:  Color.Enum(Color_BLACK).String(),
			UserID: user.ID,
		}
		s.Db.Create(&p)
		game.BlackPlayerID = sql.NullInt32{Valid: true, Int32: int32(p.ID)}
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
func (s *Server) JoinGame(context.Context, *JoinGameRequest) (*JoinGameResponse, error) { // TODO: end
	panic("TODO: unfinished")
}

// Move Moves a player piece
func (s *Server) Move(context.Context, *MoveRequest) (*MoveResponse, error) { // TODO: end
	panic("TODO: unfinished")
}
