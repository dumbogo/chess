package api

import context "context"

type Server struct {
	UnimplementedChessServiceServer
}

func (s *Server) StartGame(context.Context, *StartGameRequest) (*StartGameResponse, error) { // TODO: end
	panic("TODO: unfinished")
}

func (s *Server) JoinGame(context.Context, *JoinGameRequest) (*JoinGameResponse, error) { // TODO: end
	panic("TODO: unfinished")
}

func (s *Server) Move(context.Context, *MoveRequest) (*MoveResponse, error) { // TODO: end
	panic("TODO: unfinished")
}
