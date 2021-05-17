package engine

import (
	"math"
)

// Piece piece type
type Piece interface {
	// Name Piece name
	Identifier() PieceIdentifier
	// Color piece color
	Color() Color
	// CanMove check if piece can be moved from, to Square,
	// returns true if it's possible, even if a piece can be eaten
	CanMove(board Board, movements []Movement, from, to Square) bool

	String() string
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

type king struct {
	piece
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

func kingEatableInSquare(color Color, b Board, m []Movement, to Square) bool {
	// Loop to find if any piece can eat king in TO Square
	for _, square := range b.Squares() {
		if !square.Empty && square.Piece.Color() != color {
			if square.Piece.CanMove(b, m, square, to) {
				return true
			}
		}
	}
	return false
}

func (k *king) CanMove(b Board, m []Movement, from, to Square) bool {
	if from.Empty {
		return false
	}
	if from.Piece.Identifier() != KingIdentifier {
		return false
	}

	if !((to.Coordinates.X == from.Coordinates.X+1 || to.Coordinates.X == from.Coordinates.X-1) &&
		(to.Coordinates.Y == from.Coordinates.Y+1 || to.Coordinates.Y == from.Coordinates.Y-1) &&
		(to.Empty || !to.Empty && to.Piece.Color() != from.Piece.Color())) {
		return false
	}

	if kingEatableInSquare(k.Color(), b, m, to) {
		return false
	}

	return true
}

func (k *king) String() string {
	if k.Color() == WhiteColor {
		return "WK"
	}
	return "BK"
}

type queen struct {
	piece
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

func (q *queen) CanMove(b Board, m []Movement, from, to Square) bool {
	if from.Empty {
		return false
	}
	squares := b.Squares()
	if validRookMovement(squares, from, to) || validBishopMovement(squares, from, to) {
		return true
	}
	return false
}

func (q *queen) String() string {
	if q.Color() == WhiteColor {
		return "WQ"
	}
	return "BQ"
}

type rook struct {
	piece
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

func validRookMovement(squares Squares, from, to Square) bool {
	if from.Coordinates.X == to.Coordinates.X &&
		from.Coordinates.Y != to.Coordinates.Y {

		if !to.Empty && to.Piece.Color() == from.Piece.Color() {
			return false
		}

		// we need to make sure that, for each square between from and to, are empty
		summ := 1
		if from.Coordinates.Y > to.Coordinates.Y {
			summ = -1
		}
		for Yiterator := (int(from.Coordinates.Y) + summ); Yiterator != int(to.Coordinates.Y); Yiterator += summ {
			square := squares[CoordinateToSquareIdentifier(Coordinate{X: from.Coordinates.X, Y: uint8(Yiterator)})]
			if !square.Empty {
				return false
			}
		}
		return true
	} else if from.Coordinates.Y == to.Coordinates.Y && from.Coordinates.X != to.Coordinates.X {
		if !to.Empty && to.Piece.Color() == from.Piece.Color() {
			return false
		}

		// we need to make sure that, for each square between from and to, are empty
		summ := 1
		if from.Coordinates.X > to.Coordinates.X {
			summ = -1
		}
		for Xiterator := (int(from.Coordinates.X) + summ); Xiterator != int(to.Coordinates.X); Xiterator += summ {
			square := squares[CoordinateToSquareIdentifier(Coordinate{X: uint8(Xiterator), Y: from.Coordinates.Y})]
			if !square.Empty {
				return false
			}
		}
		return true

	}
	return false
}
func (r *rook) CanMove(b Board, m []Movement, from, to Square) bool {
	if from.Empty {
		return false
	}
	squares := b.Squares()
	return validRookMovement(squares, from, to)
}

func (r *rook) String() string {
	if r.Color() == WhiteColor {
		return "WR"
	}
	return "BR"
}

type bishop struct {
	piece
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

func validBishopMovement(squares Squares, from, to Square) bool {
	if math.Abs(float64(from.Coordinates.X)-float64(to.Coordinates.X)) != math.Abs(float64(from.Coordinates.Y)-float64(to.Coordinates.Y)) {
		return false
	}
	sumX := 1
	sumY := 1
	if from.Coordinates.X > to.Coordinates.X {
		sumX = -1
	}
	if from.Coordinates.Y > to.Coordinates.Y {
		sumY = -1
	}

	for itX, itY := int(from.Coordinates.X)+sumX, int(from.Coordinates.Y)+sumY; uint8(itX) != to.Coordinates.X && uint8(itY) != to.Coordinates.Y; {
		square := squares[CoordinateToSquareIdentifier(Coordinate{uint8(itX), uint8(itY)})]

		if !square.Empty {
			return false
		}
		itX += sumX
		itY += sumY
	}
	if !to.Empty && to.Piece.Color() == from.Piece.Color() {
		return false
	}
	return true

}
func (bi *bishop) CanMove(b Board, m []Movement, from, to Square) bool {
	//ENHANCEMENT: We can refactor and use Slope(https://en.wikipedia.org/wiki/Slope) instead
	if from.Empty {
		return false
	}
	squares := b.Squares()
	return validBishopMovement(squares, from, to)
}

func (bi *bishop) String() string {
	if bi.Color() == WhiteColor {
		return "WB"
	}
	return "BB"
}

type knight struct {
	piece
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

func (k *knight) CanMove(b Board, m []Movement, from, to Square) bool {
	if from.Empty {
		return false
	}
	if from.Piece.Identifier() != KnightIdentifier {
		return false
	}
	if !to.Empty && to.Piece.Color() == from.Piece.Color() {
		return false
	}
	// ENHANCEMENT: We can use slope, and length of vertex instead in order to recognize valid movements
	// Doing nasty if statements
	switch to.Coordinates.X {
	case from.Coordinates.X - 1:
		if to.Coordinates.Y == from.Coordinates.Y+2 || to.Coordinates.Y == from.Coordinates.Y-2 {
			return true
		}
	case from.Coordinates.X + 1:
		if to.Coordinates.Y == from.Coordinates.Y+2 || to.Coordinates.Y == from.Coordinates.Y-2 {
			return true
		}
	case from.Coordinates.X + 2:
		if to.Coordinates.Y == from.Coordinates.Y-1 || to.Coordinates.Y == from.Coordinates.Y+1 {
			return true
		}
	case from.Coordinates.X - 2:
		if to.Coordinates.Y == from.Coordinates.Y-1 || to.Coordinates.Y == from.Coordinates.Y+1 {
			return true
		}
	}
	return false
}

func (k *knight) String() string {
	if k.Color() == WhiteColor {
		return "Wk"
	}
	return "Bk"
}

type pawn struct {
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

// TODO: PawnCanMove [En passant , Handle when reachs Coordinates limits]
func (p *pawn) CanMove(board Board, movements []Movement, from, to Square) bool {
	if from.Empty {
		return false
	}
	if from.Piece.Identifier() != PawnIdentifier {
		return false
	}

	if from.Piece.Color() == WhiteColor {
		// First movement , not eating
		if from.Coordinates.Y == 1 &&
			from.Coordinates.X == to.Coordinates.X &&
			to.Coordinates.Y >= from.Coordinates.Y+2 &&
			to.Empty {
			// BUG : when there's a piece in position X+1 and pawn moves two spaces, it returns true
			return true
		}

		// Simple movement, advance one
		if to.Empty &&
			from.Coordinates.Y+1 == to.Coordinates.Y &&
			from.Coordinates.X == to.Coordinates.X {
			return true
		}

		// Eating piece
		if from.Coordinates.Y+1 == to.Coordinates.Y &&
			(from.Coordinates.X-1 == to.Coordinates.X || from.Coordinates.X+1 == to.Coordinates.X) &&
			!to.Empty &&
			to.Piece.Color() != from.Piece.Color() {
			return true
		}
	} else if from.Piece.Color() == BlackColor {
		// First movement , not eating
		if from.Coordinates.Y == 6 &&
			from.Coordinates.X == to.Coordinates.X &&
			to.Coordinates.Y >= from.Coordinates.Y-2 &&
			to.Empty {
			// BUG : when there's a piece in position X-1 and pawn moves two spaces, it returns true
			return true
		}

		// Simple movement, advance one
		if to.Empty &&
			from.Coordinates.Y-1 == to.Coordinates.Y &&
			from.Coordinates.X == to.Coordinates.X {
			return true
		}

		// Eating piece
		if from.Coordinates.Y-1 == to.Coordinates.Y &&
			(from.Coordinates.X-1 == to.Coordinates.X || from.Coordinates.X+1 == to.Coordinates.X) &&
			!to.Empty &&
			to.Piece.Color() != from.Piece.Color() {
			return true
		}
	}
	return false
}

func (p *pawn) String() string {
	if p.Color() == WhiteColor {
		return "WP"
	}
	return "BP"
}

// PieceFromPieceIdentifier returns the corresponding Piece associated to PieceIdentifier
func PieceFromPieceIdentifier(i PieceIdentifier, color Color) Piece {
	switch i {
	case PawnIdentifier:
		return NewPawn(color)
	case BishopIdentifier:
		return NewBishop(color)
	case KnightIdentifier:
		return NewKnight(color)
	case RookIdentifier:
		return NewRook(color)
	case QueenIdentifier:
		return NewQueen(color)
	case KingIdentifier:
		return NewKing(color)
	default:
		return nil
	}
}
