package engine

import (
	"strings"

	"github.com/olekukonko/tablewriter"
)

// Board chess playable board
type Board interface {
	// EatPiece removes a piece from Board location return Piece
	EatPiece(loc SquareIdentifier) Piece
	WhitePlayer() *Player
	BlackPlayer() *Player
	Squares() Squares
	String() string
}

type board struct {
	whitePlayer *Player

	blackPlayer *Player
	squares     Squares
}

// NewBoard creates a new Board to play
func NewBoard(whitePlayer, blackPlayer *Player) Board {
	board := board{
		blackPlayer: blackPlayer,
		whitePlayer: whitePlayer,
	}
	board.fillSquares()
	return &board
}

func (b *board) WhitePlayer() *Player {
	return b.whitePlayer
}

func (b *board) BlackPlayer() *Player {
	return b.blackPlayer
}

// FillSquares fill squares board with new game
func (b *board) fillSquares() {
	b.squares = Squares{
		// White pieces
		A1: Square{Empty: false, Coordinates: Coordinate{0, 0, A1}, SquareIdentifier: A1, Piece: NewRook(WhiteColor)},
		B1: Square{Empty: false, Coordinates: Coordinate{1, 0, B1}, SquareIdentifier: B1, Piece: NewKnight(WhiteColor)},
		C1: Square{Empty: false, Coordinates: Coordinate{2, 0, C1}, SquareIdentifier: C1, Piece: NewBishop(WhiteColor)},
		D1: Square{Empty: false, Coordinates: Coordinate{3, 0, D1}, SquareIdentifier: D1, Piece: NewQueen(WhiteColor)},
		E1: Square{Empty: false, Coordinates: Coordinate{4, 0, E1}, SquareIdentifier: E1, Piece: NewKing(WhiteColor)},
		F1: Square{Empty: false, Coordinates: Coordinate{5, 0, F1}, SquareIdentifier: F1, Piece: NewBishop(WhiteColor)},
		G1: Square{Empty: false, Coordinates: Coordinate{6, 0, G1}, SquareIdentifier: G1, Piece: NewKnight(WhiteColor)},
		H1: Square{Empty: false, Coordinates: Coordinate{7, 0, H1}, SquareIdentifier: H1, Piece: NewRook(WhiteColor)},
		A2: Square{Empty: false, Coordinates: Coordinate{0, 1, A2}, SquareIdentifier: A2, Piece: NewPawn(WhiteColor)},
		B2: Square{Empty: false, Coordinates: Coordinate{1, 1, B2}, SquareIdentifier: B2, Piece: NewPawn(WhiteColor)},
		C2: Square{Empty: false, Coordinates: Coordinate{2, 1, C2}, SquareIdentifier: C2, Piece: NewPawn(WhiteColor)},
		D2: Square{Empty: false, Coordinates: Coordinate{3, 1, D2}, SquareIdentifier: D2, Piece: NewPawn(WhiteColor)},
		E2: Square{Empty: false, Coordinates: Coordinate{4, 1, E2}, SquareIdentifier: E2, Piece: NewPawn(WhiteColor)},
		F2: Square{Empty: false, Coordinates: Coordinate{5, 1, F2}, SquareIdentifier: F2, Piece: NewPawn(WhiteColor)},
		G2: Square{Empty: false, Coordinates: Coordinate{6, 1, G2}, SquareIdentifier: G2, Piece: NewPawn(WhiteColor)},
		H2: Square{Empty: false, Coordinates: Coordinate{7, 1, H2}, SquareIdentifier: H2, Piece: NewPawn(WhiteColor)},

		// Black pieces
		A8: Square{Empty: false, Coordinates: Coordinate{0, 7, A8}, SquareIdentifier: A8, Piece: NewRook(BlackColor)},
		B8: Square{Empty: false, Coordinates: Coordinate{1, 7, B8}, SquareIdentifier: B8, Piece: NewKnight(BlackColor)},
		C8: Square{Empty: false, Coordinates: Coordinate{2, 7, C8}, SquareIdentifier: C8, Piece: NewBishop(BlackColor)},
		D8: Square{Empty: false, Coordinates: Coordinate{3, 7, D8}, SquareIdentifier: D8, Piece: NewQueen(BlackColor)},
		E8: Square{Empty: false, Coordinates: Coordinate{4, 7, E8}, SquareIdentifier: E8, Piece: NewKing(BlackColor)},
		F8: Square{Empty: false, Coordinates: Coordinate{5, 7, F8}, SquareIdentifier: F8, Piece: NewBishop(BlackColor)},
		G8: Square{Empty: false, Coordinates: Coordinate{6, 7, G8}, SquareIdentifier: G8, Piece: NewKnight(BlackColor)},
		H8: Square{Empty: false, Coordinates: Coordinate{7, 7, H8}, SquareIdentifier: H8, Piece: NewRook(BlackColor)},
		A7: Square{Empty: false, Coordinates: Coordinate{0, 6, A7}, SquareIdentifier: A7, Piece: NewPawn(BlackColor)},
		B7: Square{Empty: false, Coordinates: Coordinate{1, 6, B7}, SquareIdentifier: B7, Piece: NewPawn(BlackColor)},
		C7: Square{Empty: false, Coordinates: Coordinate{2, 6, C7}, SquareIdentifier: C7, Piece: NewPawn(BlackColor)},
		D7: Square{Empty: false, Coordinates: Coordinate{3, 6, D7}, SquareIdentifier: D7, Piece: NewPawn(BlackColor)},
		E7: Square{Empty: false, Coordinates: Coordinate{4, 6, E7}, SquareIdentifier: E7, Piece: NewPawn(BlackColor)},
		F7: Square{Empty: false, Coordinates: Coordinate{5, 6, F7}, SquareIdentifier: F7, Piece: NewPawn(BlackColor)},
		G7: Square{Empty: false, Coordinates: Coordinate{6, 6, G7}, SquareIdentifier: G7, Piece: NewPawn(BlackColor)},
		H7: Square{Empty: false, Coordinates: Coordinate{7, 6, H7}, SquareIdentifier: H7, Piece: NewPawn(BlackColor)},

		// Empty pieces
		A3: Square{Empty: true, Coordinates: Coordinate{0, 2, A3}, SquareIdentifier: A3},
		B3: Square{Empty: true, Coordinates: Coordinate{1, 2, B3}, SquareIdentifier: B3},
		C3: Square{Empty: true, Coordinates: Coordinate{2, 2, C3}, SquareIdentifier: C3},
		D3: Square{Empty: true, Coordinates: Coordinate{3, 2, D3}, SquareIdentifier: D3},
		E3: Square{Empty: true, Coordinates: Coordinate{4, 2, E3}, SquareIdentifier: E3},
		F3: Square{Empty: true, Coordinates: Coordinate{5, 2, F3}, SquareIdentifier: F3},
		G3: Square{Empty: true, Coordinates: Coordinate{6, 2, G3}, SquareIdentifier: G3},
		H3: Square{Empty: true, Coordinates: Coordinate{7, 2, H3}, SquareIdentifier: H3},

		A4: Square{Empty: true, Coordinates: Coordinate{0, 3, A4}, SquareIdentifier: A4},
		B4: Square{Empty: true, Coordinates: Coordinate{1, 3, B4}, SquareIdentifier: B4},
		C4: Square{Empty: true, Coordinates: Coordinate{2, 3, C4}, SquareIdentifier: C4},
		D4: Square{Empty: true, Coordinates: Coordinate{3, 3, D4}, SquareIdentifier: D4},
		E4: Square{Empty: true, Coordinates: Coordinate{4, 3, E4}, SquareIdentifier: E4},
		F4: Square{Empty: true, Coordinates: Coordinate{5, 3, F4}, SquareIdentifier: F4},
		G4: Square{Empty: true, Coordinates: Coordinate{6, 3, G4}, SquareIdentifier: G4},
		H4: Square{Empty: true, Coordinates: Coordinate{7, 3, H4}, SquareIdentifier: H4},

		A5: Square{Empty: true, Coordinates: Coordinate{0, 4, A5}, SquareIdentifier: A5},
		B5: Square{Empty: true, Coordinates: Coordinate{1, 4, B5}, SquareIdentifier: B5},
		C5: Square{Empty: true, Coordinates: Coordinate{2, 4, C5}, SquareIdentifier: C5},
		D5: Square{Empty: true, Coordinates: Coordinate{3, 4, D5}, SquareIdentifier: D5},
		E5: Square{Empty: true, Coordinates: Coordinate{4, 4, E5}, SquareIdentifier: E5},
		F5: Square{Empty: true, Coordinates: Coordinate{5, 4, F5}, SquareIdentifier: F5},
		G5: Square{Empty: true, Coordinates: Coordinate{6, 4, G5}, SquareIdentifier: G5},
		H5: Square{Empty: true, Coordinates: Coordinate{7, 4, H5}, SquareIdentifier: H5},

		A6: Square{Empty: true, Coordinates: Coordinate{0, 5, A6}, SquareIdentifier: A6},
		B6: Square{Empty: true, Coordinates: Coordinate{1, 5, B6}, SquareIdentifier: B6},
		C6: Square{Empty: true, Coordinates: Coordinate{2, 5, C6}, SquareIdentifier: C6},
		D6: Square{Empty: true, Coordinates: Coordinate{3, 5, D6}, SquareIdentifier: D6},
		E6: Square{Empty: true, Coordinates: Coordinate{4, 5, E6}, SquareIdentifier: E6},
		F6: Square{Empty: true, Coordinates: Coordinate{5, 5, F6}, SquareIdentifier: F6},
		G6: Square{Empty: true, Coordinates: Coordinate{6, 5, G6}, SquareIdentifier: G6},
		H6: Square{Empty: true, Coordinates: Coordinate{7, 5, H6}, SquareIdentifier: H6},
	}
}

func (b *board) EatPiece(loc SquareIdentifier) Piece {
	square := b.squares[loc]
	piece := square.Piece
	square.Piece = nil
	square.Empty = true
	b.squares[loc] = square
	return piece
}

func (b *board) Squares() Squares {
	return b.squares
}

func (b *board) String() string {
	var builder strings.Builder

	a1, b1, c1, d1, e1, f1, g1, h1 := b.squares[A1], b.squares[B1], b.squares[C1], b.squares[D1], b.squares[E1], b.squares[F1], b.squares[G1], b.squares[H1]
	a2, b2, c2, d2, e2, f2, g2, h2 := b.squares[A2], b.squares[B2], b.squares[C2], b.squares[D2], b.squares[E2], b.squares[F2], b.squares[G2], b.squares[H2]
	a3, b3, c3, d3, e3, f3, g3, h3 := b.squares[A3], b.squares[B3], b.squares[C3], b.squares[D3], b.squares[E3], b.squares[F3], b.squares[G3], b.squares[H3]
	a4, b4, c4, d4, e4, f4, g4, h4 := b.squares[A4], b.squares[B4], b.squares[C4], b.squares[D4], b.squares[E4], b.squares[F4], b.squares[G4], b.squares[H4]
	a5, b5, c5, d5, e5, f5, g5, h5 := b.squares[A5], b.squares[B5], b.squares[C5], b.squares[D5], b.squares[E5], b.squares[F5], b.squares[G5], b.squares[H5]
	a6, b6, c6, d6, e6, f6, g6, h6 := b.squares[A6], b.squares[B6], b.squares[C6], b.squares[D6], b.squares[E6], b.squares[F6], b.squares[G6], b.squares[H6]
	a7, b7, c7, d7, e7, f7, g7, h7 := b.squares[A7], b.squares[B7], b.squares[C7], b.squares[D7], b.squares[E7], b.squares[F7], b.squares[G7], b.squares[H7]
	a8, b8, c8, d8, e8, f8, g8, h8 := b.squares[A8], b.squares[B8], b.squares[C8], b.squares[D8], b.squares[E8], b.squares[F8], b.squares[G8], b.squares[H8]

	data := [][]string{
		{"8", a8.String(), b8.String(), c8.String(), d8.String(), e8.String(), f8.String(), g8.String(), h8.String()},
		{"7", a7.String(), b7.String(), c7.String(), d7.String(), e7.String(), f7.String(), g7.String(), h7.String()},
		{"6", a6.String(), b6.String(), c6.String(), d6.String(), e6.String(), f6.String(), g6.String(), h6.String()},
		{"5", a5.String(), b5.String(), c5.String(), d5.String(), e5.String(), f5.String(), g5.String(), h5.String()},
		{"4", a4.String(), b4.String(), c4.String(), d4.String(), e4.String(), f4.String(), g4.String(), h4.String()},
		{"3", a3.String(), b3.String(), c3.String(), d3.String(), e3.String(), f3.String(), g3.String(), h3.String()},
		{"2", a2.String(), b2.String(), c2.String(), d2.String(), e2.String(), f2.String(), g2.String(), h2.String()},
		{"1", a1.String(), b1.String(), c1.String(), d1.String(), e1.String(), f1.String(), g1.String(), h1.String()},
		{"", "a", "b", "c", "d", "e", "f", "g", "h"},
	}

	table := tablewriter.NewWriter(&builder)

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
	return builder.String()
}

// Squares board squares, total of 8*8
type Squares map[SquareIdentifier]Square

// Square type val
type Square struct {
	Empty            bool
	Piece            Piece
	Coordinates      Coordinate
	SquareIdentifier SquareIdentifier
}

func (s *Square) String() string {
	if s.Empty {
		return ""
	}
	return s.Piece.String()
}

// MAXY maximum coordinate value
const MAXY = 7

// MAXX maximum coordinate value
const MAXX = 7

// Coordinate a coordinate within the Board, starting at 0,0 = A1, 8,8 = H8
type Coordinate struct {
	X                uint8
	Y                uint8
	SquareIdentifier SquareIdentifier
}
