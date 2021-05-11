package engine

import (
	"errors"
)

// TODO: Pawn Promotion
// TODO: check if En passant to delete Pawn
// TODO: check Castling(short & large)

// Player gammer
type Player struct {
	Name  string
	Color Color
}

// Movement a movement made
type Movement struct {
	Player     Player
	PieceMoved Piece
	PieceEaten Piece
	From       SquareIdentifier
	To         SquareIdentifier
}

// Game playable game
type Game interface {
	// Turn returns player turn
	Turn() Player
	// Move moves a piece in the Board, returns true if moved
	Move(player Player, form, to SquareIdentifier) (bool, error)
	// IsCheckBy returns true if Player makes check
	IsCheckBy(Player) bool
	// IsCheckmateBy returns true if Player makes checkmate
	IsCheckmateBy(Player) bool
	// Board get board
	Board() Board
	// Movements get all historic movements
	Movements() []Movement

	String() string

	WhitePieces() map[PieceIdentifier]uint8
	BlackPieces() map[PieceIdentifier]uint8
}

type game struct {
	name  string
	board Board
	turn  Player
	white Player
	black Player

	whitePieces map[PieceIdentifier]uint8
	blackPieces map[PieceIdentifier]uint8
	movements   []Movement
}

// NewGame creates new Game
func NewGame(name string, black, white Player) (Game, error) {
	if white.Color != WhiteColor || black.Color != BlackColor {
		return nil, errors.New("must define black and white players")
	}

	blackPieces := map[PieceIdentifier]uint8{
		RookIdentifier:   2,
		KnightIdentifier: 2,
		BishopIdentifier: 2,
		QueenIdentifier:  1,
		KingIdentifier:   1,
		PawnIdentifier:   8,
	}
	whitePieces := map[PieceIdentifier]uint8{
		RookIdentifier:   2,
		KnightIdentifier: 2,
		BishopIdentifier: 2,
		QueenIdentifier:  1,
		KingIdentifier:   1,
		PawnIdentifier:   8,
	}
	return &game{
		name:        name,
		board:       NewBoard(&white, &black),
		turn:        white,
		white:       white,
		black:       black,
		blackPieces: blackPieces,
		whitePieces: whitePieces,
	}, nil
}

// LoadGame loads in game being played
func LoadGame(
	name string,
	board Board,
	turn Player,
	white, black Player,
	whitePieces map[PieceIdentifier]uint8,
	blackPieces map[PieceIdentifier]uint8,
	movements []Movement,
) (Game, error) {
	return &game{
		name:        name,
		board:       board,
		turn:        turn,
		white:       white,
		black:       black,
		whitePieces: whitePieces,
		blackPieces: blackPieces,
		movements:   movements,
	}, nil
}

func (g *game) Turn() Player {
	return g.turn
}

func (g *game) WhitePieces() map[PieceIdentifier]uint8 {
	return g.whitePieces
}

func (g *game) BlackPieces() map[PieceIdentifier]uint8 {
	return g.blackPieces
}

func (g *game) Move(player Player, from, to SquareIdentifier) (bool, error) {
	squareFrom := g.board.Squares()[from]
	squareTo := g.board.Squares()[to]
	if squareFrom.Empty {
		return false, nil
	}

	pieceToMove := squareFrom.Piece
	if pieceToMove.Color() != g.Turn().Color {
		return false, nil
	}

	canMove := squareFrom.Piece.CanMove(g.board, g.Movements(), squareFrom, squareTo)
	if !canMove {
		return false, nil
	}

	var pieceEaten Piece
	if !squareTo.Empty {
		pieceEaten = squareTo.Piece
		g.board.EatPiece(to)
		g.removePiecePlayer(pieceEaten)
	}

	squareTo.Piece = pieceToMove
	squareTo.Empty = false
	squareFrom.Piece = nil
	squareFrom.Empty = true
	g.board.Squares()[to] = squareTo
	g.board.Squares()[from] = squareFrom
	g.changeTurn()
	g.movements = append(g.movements, Movement{
		Player:     player,
		PieceMoved: pieceToMove,
		PieceEaten: pieceEaten,
		From:       from,
		To:         to,
	})
	return true, nil
}

func (g *game) Movements() []Movement {
	return g.movements
}

func (g *game) IsCheckBy(player Player) bool {
	var kingSquare Square
	playerColor := player.Color
	if playerColor == WhiteColor {
		kingSquare = getKingSquare(g.Board(), BlackColor)
	} else {
		kingSquare = getKingSquare(g.Board(), WhiteColor)
	}
	if kingEatableInSquare(kingSquare.Piece.Color(), g.Board(), g.Movements(), kingSquare) {
		return true
	}
	return false
}

func (g *game) IsCheckmateBy(player Player) bool {
	if !g.IsCheckBy(player) {
		return false
	}
	var kingSquare Square
	playerColor := player.Color
	if playerColor == WhiteColor {
		kingSquare = getKingSquare(g.Board(), BlackColor)
	} else {
		kingSquare = getKingSquare(g.Board(), WhiteColor)
	}
	// iterate all possible movements and checks if king can move to at least one direction
	for x := kingSquare.Coordinates.X - 1; x <= kingSquare.Coordinates.X+1; x++ {
		for y := kingSquare.Coordinates.Y - 1; y <= kingSquare.Coordinates.Y+1; y++ {
			if (x > MAXX || y > MAXY) || x == kingSquare.Coordinates.X && y == kingSquare.Coordinates.Y {
				continue
			}
			if kingSquare.Piece.CanMove(g.Board(), g.Movements(), kingSquare, g.Board().Squares()[CoordinateToSquareIdentifier(Coordinate{x, y})]) {
				return false
			}
		}
	}
	return true
}

func (g *game) Board() Board {
	return g.board
}

// String returns ASCII representation of the game
func (g *game) String() string {
	return g.board.String()
}

func (g *game) changeTurn() {
	if g.Turn().Color == WhiteColor {
		g.turn = g.black
	} else {
		g.turn = g.white
	}
}

func (g *game) removePiecePlayer(p Piece) {
	color := p.Color()
	switch color {
	case BlackColor:
		g.blackPieces[p.Identifier()]--
	case WhiteColor:
		g.whitePieces[p.Identifier()]--
	}
}

func getKingSquare(board Board, color Color) Square {
	var kingSquare Square
	for _, square := range board.Squares() {
		if !square.Empty && square.Piece.Identifier() == KingIdentifier && square.Piece.Color() == color {
			kingSquare = square
			break
		}
	}
	return kingSquare
}
