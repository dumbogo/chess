package api

import (
	context "context"

	"gorm.io/gorm"
)

// Server grpc server interface implementation
type Server struct {
	UnimplementedChessServiceServer
	Db *gorm.DB
}

// StartGame starts a new game
func (s *Server) StartGame(context.Context, *StartGameRequest) (*StartGameResponse, error) { // TODO: end
	panic("TODO: unfinished")
}

// JoinGame joins a game
func (s *Server) JoinGame(context.Context, *JoinGameRequest) (*JoinGameResponse, error) { // TODO: end
	panic("TODO: unfinished")
}

// Move Moves a player piece
func (s *Server) Move(context.Context, *MoveRequest) (*MoveResponse, error) { // TODO: end
	panic("TODO: unfinished")
}
