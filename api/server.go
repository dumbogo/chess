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

func newGameWithoutPlayers(name string) Game {
	engineSquares := engine.PristineSquares()
	squares := make(Squares)
	for i, v := range engineSquares {
		sq := newBasicSquare(v)
		if !sq.Empty {
			sq.Piece = Piece{
				PieceIdentifier: v.Piece.Identifier(),
				Color:           v.Piece.Color(),
			}
		}
		squares[i] = sq
	}
	return Game{
		Name: name,
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
		BoardSquares: squares,
	}
}

func newBasicSquare(e engine.Square) Square {
	return Square{e, Piece{}}
}

// StartGame starts a new game
func (s *Server) StartGame(ctx context.Context, startGameRequest *StartGameRequest) (*StartGameResponse, error) {
	user, e := getUserFromCtx(ctx)
	if e != nil {
		return nil, e
	}
	game := newGameWithoutPlayers(startGameRequest.GetName())
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

// JoinGame joins a game, depending of the white or black player space left, it will be assigned to the user joining
func (s *Server) JoinGame(ctx context.Context, r *JoinGameRequest) (*JoinGameResponse, error) {
	user, e := getUserFromCtx(ctx)
	if e != nil {
		return nil, e
	}
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
func (s *Server) Move(ctx context.Context, r *MoveRequest) (*MoveResponse, error) {
	user, e := getUserFromCtx(ctx)
	if e != nil {
		return nil, e
	}
	game := Game{}
	tx := DBConn.Where("uuid=?", r.GetUuid()).Joins("join players on players.id = white_player_id").First(&game)
	if tx.Error != nil {
		return nil, tx.Error
	}
	whitePlayerDb := Player{}
	tx = DBConn.Where("id=?", game.WhitePlayerID.Int32).First(&whitePlayerDb)
	if tx.Error != nil {
		return nil, tx.Error
	}
	blackPlayerDb := Player{}
	tx = DBConn.Where("id=?", game.BlackPlayerID.Int32).First(&blackPlayerDb)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if uint(game.Turn) == whitePlayerDb.ID {
		if user.ID != whitePlayerDb.UserID {
			return nil, errors.New("not your turn")
		}
	}

	nextTurn := uint(whitePlayerDb.ID)
	turnPlayer := engine.Player{
		Color: engine.WhiteColor,
	}
	if uint(game.Turn) == whitePlayerDb.ID {
		nextTurn = uint(blackPlayerDb.ID)
		turnPlayer.Color = engine.WhiteColor
	} else {
		turnPlayer.Color = engine.BlackColor
	}
	whitePlayer := engine.Player{Color: engine.WhiteColor}
	blackPlayer := engine.Player{Color: engine.BlackColor}

	squares := engine.Squares{}
	for i, v := range game.BoardSquares {
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
	board := engine.LoadBoard(&whitePlayer, &blackPlayer, squares)
	whitePieces := engine.PiecesList{}
	for i, v := range game.WhitePieces {
		whitePieces[engine.PieceIdentifier(i)] = v
	}
	blackPieces := engine.PiecesList{}
	for i, v := range game.BlackPieces {
		blackPieces[engine.PieceIdentifier(i)] = v
	}

	gameEngine, e := engine.LoadGame(
		game.Name,
		board,
		turnPlayer,
		whitePlayer,
		blackPlayer,
		whitePieces,
		blackPieces,
		make([]engine.Movement, 0), // Movements, leaving emtpy ATM
	)
	if e != nil {
		return nil, e
	}

	from, ok := engine.StringToSquareIdentifier(r.GetFromSquare())
	if !ok {
		return nil, errors.New("Someerr")
	}
	to, ok := engine.StringToSquareIdentifier(r.GetToSquare())
	if !ok {
		return nil, errors.New("Someerr")
	}
	ok, e = gameEngine.Move(
		turnPlayer,
		from,
		to,
	)
	if !ok {
		return nil, e
	}

	// Update values on gameDB
	game.Turn = nextTurn
	newBlackPieces := pieces{}
	for i, v := range gameEngine.BlackPieces() {
		newBlackPieces[uint8(i)] = v
	}
	game.BlackPieces = newBlackPieces

	newWhitePieces := pieces{}
	for i, v := range gameEngine.WhitePieces() {
		newWhitePieces[uint8(i)] = v
	}
	game.WhitePieces = newWhitePieces

	newSquares := Squares{}
	for i, v := range gameEngine.Board().Squares() {
		sq := newBasicSquare(v)
		if !sq.Empty {
			sq.Piece = Piece{
				PieceIdentifier: v.Piece.Identifier(),
				Color:           v.Piece.Color(),
			}
		}
		newSquares[i] = sq
	}
	game.BoardSquares = newSquares
	tx = DBConn.Save(&game)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &MoveResponse{
		Board: gameEngine.Board().String(),
	}, nil
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

func getUserFromCtx(ctx context.Context) (User, error) {
	t, e := GetAccessTokenFromCtx(ctx)
	if e != nil {
		return User{}, e
	}
	user := GetUserFromAccessToken(t)
	return user, nil
}
