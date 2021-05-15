package api

import (
	context "context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/dumbogo/chess/engine"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	status "google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

// Server grpc server interface implementation
type Server struct {
	UnimplementedChessServiceServer
	Db *gorm.DB
}

// StartGame starts a new game
func (s *Server) StartGame(ctx context.Context, startGameRequest *StartGameRequest) (*StartGameResponse, error) {
	accessToken, err := GetAccessTokenFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	user := GetUserFromAccessToken(accessToken)
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
	accessToken, err := GetAccessTokenFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	user := GetUserFromAccessToken(accessToken)

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

// EnsureValidToken ensures a valid token exists within a request's metadata. If
// the token is missing or invalid, the interceptor blocks execution of the
// handler and returns an error. Otherwise, the interceptor invokes the unary
// handler.
func EnsureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	// The keys within metadata.MD are normalized to lowercase.
	// See: https://godoc.org/google.golang.org/grpc/metadata#New
	if !valid(md["authorization"]) {
		return nil, errInvalidToken
	}
	// Continue execution of handler after ensuring a valid token.
	return handler(ctx, req)
}

// valid validates the authorization. Returns false if neither user, nor auth token found and token expired
func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	accessToken := strings.TrimPrefix(authorization[0], "Bearer ")
	user := GetUserFromAccessToken(accessToken)
	if user.Email == "" || user.ExpiresAt.Valid && time.Now().After(user.ExpiresAt.Time) {
		return false
	}
	return true
}

// GetAccessTokenFromCtx returns access token from ctx object metadata
func GetAccessTokenFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errMissingMetadata
	}
	authorization := md["authorization"]
	accessToken := strings.TrimPrefix(authorization[0], "Bearer ")
	return accessToken, nil
}
