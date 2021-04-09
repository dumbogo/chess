package engine

import (
	"errors"
)

// Piece piece type
type Piece interface {
	// Name Piece name
	Identifier() PieceIdentifier
	// Color piece color
	Color() Color
	// CanMove check if piece can be moved from, to Square,
	// returns true if it's possible
	CanMove(board Board, from, to Square) bool
	PossibleMovements(board Board, from Square) []Square
}

type piece struct {
	pieceIdentifier PieceIdentifier
	color           Color
}

// Name interface Piece
func (p *piece) Identifier() PieceIdentifier {
	return p.pieceIdentifier
}

func (p *piece) Color() Color {
	return p.color
}

func (p *piece) CanMove(b Board, from, to Square) bool {
	panic(errors.New("Must implement interface method"))
}

func (p *piece) PossibleMovements(Board, Square) []Square {
	panic(errors.New("Must implement interface method"))
}

type pawn struct {
	piece
}

type bishop struct {
	piece
}

type knight struct {
	piece
}

type rook struct {
	piece
}

type queen struct {
	piece
}

type king struct {
	piece
}

// NewPawn creates Pawn piece
func NewPawn(color Color) Piece {
	p := piece{
		pieceIdentifier: PawnIdentifier,
		color:           color,
	}
	return &pawn{
		piece: p,
	}
}

// NewBishop creates Bishop piece
func NewBishop(color Color) Piece {
	p := piece{
		pieceIdentifier: BishopIdentifier,
		color:           color,
	}
	return &bishop{
		piece: p,
	}
}

// NewRook creates Rook piece
func NewRook(color Color) Piece {
	p := piece{
		pieceIdentifier: RookIdentifier,
		color:           color,
	}
	return &rook{
		piece: p,
	}
}

// NewKnight creates knight piece
func NewKnight(color Color) Piece {
	p := piece{
		pieceIdentifier: KnightIdentifier,
		color:           color,
	}
	return &knight{
		piece: p,
	}
}

// NewQueen creates Queen piece
func NewQueen(color Color) Piece {
	p := piece{
		pieceIdentifier: QueenIdentifier,
		color:           color,
	}
	return &queen{
		piece: p,
	}
}

// NewKing creates King piece
func NewKing(color Color) Piece {
	p := piece{
		pieceIdentifier: KingIdentifier,
		color:           color,
	}
	return &king{
		piece: p,
	}
}

// TODO: WIP
func (p *pawn) CanMove(board Board, from, to Square) bool {
	// check if first movement by color
	if from.Piece.Color() == BlackColor {
		if from.Coordinates.Y == 6 {
			// Here, I know It's the first movement
		}
	}
	return true
}
