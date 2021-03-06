package api

import (
	context "context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dumbogo/chess/engine"
	"github.com/dumbogo/chess/messagebroker"
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

// MessageBroker a pub/sub messagre broker system to send updates on watchers
var MessageBroker messagebroker.MessageBroker

// Server grpc server interface implementation
type Server struct {
	UnimplementedChessServiceServer
	Db *gorm.DB
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

// Watch receive live updates on the current game loaded
func (s *Server) Watch(r *WatchRequest, stream ChessService_WatchServer) error {
	// TODO: send a message when is either check or checkmate
	if MessageBroker == nil {
		log.Printf("Message broker not initialized, prompts error")
		return errors.New("internal server error")
	}
	chMsgs, err := MessageBroker.Subscribe(r.GetUuid())
	if err != nil {
		return err
	}

	// Send for the first time, the board
	gameDb := Game{}
	tx := DBConn.Where("uuid=?", r.GetUuid()).Joins("join players on players.id = white_player_id").First(&gameDb)
	if tx.Error != nil {
		return tx.Error
	}

	gameEngine, err := loadEngineGameFromDbValues(gameDb, engine.Player{})
	if err != nil {
		return err
	}

	turnColor := "white"
	if gameDb.Turn == uint(gameDb.BlackPlayerID.Int32) {
		turnColor = "black"
	}

	if err := stream.Send(&WatchResponse{
		Status: fmt.Sprintf("%s turn", turnColor),
		Turn:   gameEngine.Turn().Name,
		Board:  gameEngine.Board().String(),
	}); err != nil {
		return err
	}

	// Wait on updates for new movements
	for msg := range chMsgs {
		payload := payloadUpdateGame{}
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return err
		}
		if err := stream.Send(&WatchResponse{
			Status: payload.Status,
			Turn:   payload.Turn,
			Board:  payload.Board,
		}); err != nil {
			return err
		}
	}
	return nil
}

// Move Moves a player piece
func (s *Server) Move(ctx context.Context, r *MoveRequest) (*MoveResponse, error) {
	log.Printf("Moving piece, req: %+v\n", r)
	user, e := getUserFromCtx(ctx)
	if e != nil {
		return nil, e
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Extract game and players from db
	// TODO: gorm joins func is not working as expected, review
	gameDb := Game{}
	tx := DBConn.Where("uuid=?", r.GetUuid()).Joins("join players on players.id = white_player_id").First(&gameDb)
	if tx.Error != nil {
		return nil, tx.Error
	}

	whitePlayerDb := Player{}
	tx = DBConn.Where("id=?", gameDb.WhitePlayerID.Int32).First(&whitePlayerDb)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("white player is missing")
		}
		return nil, tx.Error
	}

	blackPlayerDb := Player{}
	tx = DBConn.Where("id=?", gameDb.BlackPlayerID.Int32).First(&blackPlayerDb)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("black player is missing")
		}
		return nil, tx.Error
	}

	log.Printf("white player: %+v\n", whitePlayerDb)
	log.Printf("black player: %+v\n", blackPlayerDb)

	log.Printf("User : %+v\n", user)

	// validate turn color and if user is any player
	if gameDb.Turn == whitePlayerDb.ID {
		if user.ID != whitePlayerDb.UserID {
			return nil, errors.New("Not your turn")
		}
	} else if gameDb.Turn == blackPlayerDb.ID {
		if user.ID != blackPlayerDb.UserID {
			return nil, errors.New("Not your turn")
		}
	}

	nextTurn := uint(whitePlayerDb.ID)
	turnPlayer := engine.Player{
		Color: engine.WhiteColor,
	}
	if uint(gameDb.Turn) == whitePlayerDb.ID {
		nextTurn = uint(blackPlayerDb.ID)
		turnPlayer.Color = engine.WhiteColor
	} else {
		turnPlayer.Color = engine.BlackColor
	}

	gameEngine, err := loadEngineGameFromDbValues(gameDb, turnPlayer)
	if err != nil {
		return nil, err
	}

	from, ok := engine.StringToSquareIdentifier(strings.ToUpper(r.GetFromSquare()))
	if !ok {
		return nil, errors.New("Invalid \"from\" square identifier")
	}

	// Review if from piece is from player
	squareFrom := gameEngine.Board().Squares()[from]
	if !squareFrom.Empty && squareFrom.Piece.Color() != turnPlayer.Color {
		return nil, errors.New("Not your piece color")
	}

	to, ok := engine.StringToSquareIdentifier(strings.ToUpper(r.GetToSquare()))
	if !ok {
		return nil, errors.New("Invalid \"to\" square identifier")
	}

	if ok, e = gameEngine.Move(turnPlayer, from, to); !ok {
		return nil, e
	}

	gameDb.Turn = nextTurn
	if err := updateGameValuesFromGameEngine(&gameDb, gameEngine); err != nil {
		return nil, err
	}

	// TODO: Send response checkmate or check when it happens
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

func loadEngineGameFromDbValues(gameDb Game, turn engine.Player) (engine.Game, error) {
	whitePlayer := engine.Player{Color: engine.WhiteColor}
	blackPlayer := engine.Player{Color: engine.BlackColor}
	board := engine.LoadBoard(&whitePlayer, &blackPlayer, squaresToEngineSquares(gameDb.BoardSquares))
	whitePieces := engine.PiecesList{}
	for i, v := range gameDb.WhitePieces {
		whitePieces[engine.PieceIdentifier(i)] = v
	}
	blackPieces := engine.PiecesList{}
	for i, v := range gameDb.BlackPieces {
		blackPieces[engine.PieceIdentifier(i)] = v
	}

	return engine.LoadGame(
		gameDb.Name,
		board,
		turn,
		whitePlayer,
		blackPlayer,
		whitePieces,
		blackPieces,
		make([]engine.Movement, 0), // Movements, leaving empty ATM
	)
}

func updateGameValuesFromGameEngine(gameDb *Game, gameEngine engine.Game) error {
	newBlackPiecesDb := pieces{}
	for i, v := range gameEngine.BlackPieces() {
		newBlackPiecesDb[uint8(i)] = v
	}
	gameDb.BlackPieces = newBlackPiecesDb
	newWhitePiecesDb := pieces{}
	for i, v := range gameEngine.WhitePieces() {
		newWhitePiecesDb[uint8(i)] = v
	}
	gameDb.WhitePieces = newWhitePiecesDb
	gameDb.BoardSquares = engineSquaresToSquares(gameEngine.Board().Squares())
	if tx := DBConn.Save(&gameDb); tx.Error != nil {
		return tx.Error
	}
	return nil
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

func getUserFromCtx(ctx context.Context) (*User, error) {
	t, e := GetAccessTokenFromCtx(ctx)
	if e != nil {
		return nil, e
	}

	return GetUserFromAccessToken(t), nil
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
