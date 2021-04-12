package engine

import (
	"errors"
)

// TODO: list
// Pawn Promotion

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
	// IsCheckBy returns Player making check
	IsCheckBy(Player) bool
	// IsCheckmateBy returns Player making checkmate
	IsCheckmateBy(Player) bool
	// Board get board
	Board() Board
	// Movements get all historic movements
	Movements() []Movement

	String() string
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
		return nil, errors.New("Must define black and white players")
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
func (g *game) Turn() Player {
	return g.turn
}

func (g *game) Move(player Player, from, to SquareIdentifier) (bool, error) {
	squareFrom := g.board.Squares()[from]
	squareTo := g.board.Squares()[to]
	if squareFrom.Empty == true {
		return false, nil
	}

	pieceToMove := squareFrom.Piece
	if pieceToMove.Color() != g.Turn().Color {
		return false, nil
	}
	canMove := squareFrom.Piece.CanMove(g.board, squareFrom, squareTo)
	if !canMove {
		return false, nil
	}

	// TODO: check if En passant to delete Pawn
	// TODO: check Castling(short & large)
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

func (g *game) IsCheckBy(p Player) bool {
	// TODO: End
	return false
}

func (g *game) IsCheckmateBy(p Player) bool {
	// TODO: end
	return false
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
