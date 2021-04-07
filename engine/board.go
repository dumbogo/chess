package engine

// Board chess playable board
type Board interface {
	FillSquares()
	EatPiece(loc SquareIdentifier) Piece
	WhitePlayer() *Player
	BlackPlayer() *Player
	Squares() Squares
}

type board struct {
	whitePlayer *Player
	whitePieces map[PieceIdentifier]uint8

	blackPlayer *Player
	blackPieces map[PieceIdentifier]uint8
	squares     Squares
}

// NewBoard creates a new Board to play
func NewBoard(whitePlayer, blackPlayer *Player) Board {
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
	board := board{
		blackPlayer: blackPlayer,
		blackPieces: blackPieces,
		whitePlayer: whitePlayer,
		whitePieces: whitePieces,
	}
	board.FillSquares()
	return &board
}

func (b *board) WhitePlayer() *Player {
	return b.whitePlayer
}

func (b *board) BlackPlayer() *Player {
	return b.blackPlayer
}

// FillSquares fill squares board with new game
func (b *board) FillSquares() {
	b.squares = Squares{
		// White pieces
		A1: Square{Empty: false, Coordinates: Coordinate{0, 0}, SquareIdentifier: A1, Piece: NewRook(WhiteColor)},
		B1: Square{Empty: false, Coordinates: Coordinate{0, 1}, SquareIdentifier: B1, Piece: NewKnight(WhiteColor)},
		C1: Square{Empty: false, Coordinates: Coordinate{0, 2}, SquareIdentifier: C1, Piece: NewBishop(WhiteColor)},
		D1: Square{Empty: false, Coordinates: Coordinate{0, 3}, SquareIdentifier: D1, Piece: NewQueen(WhiteColor)},
		E1: Square{Empty: false, Coordinates: Coordinate{0, 4}, SquareIdentifier: E1, Piece: NewKing(WhiteColor)},
		F1: Square{Empty: false, Coordinates: Coordinate{0, 5}, SquareIdentifier: F1, Piece: NewBishop(WhiteColor)},
		G1: Square{Empty: false, Coordinates: Coordinate{0, 6}, SquareIdentifier: G1, Piece: NewKnight(WhiteColor)},
		H1: Square{Empty: false, Coordinates: Coordinate{0, 7}, SquareIdentifier: H1, Piece: NewRook(WhiteColor)},
		A2: Square{Empty: false, Coordinates: Coordinate{1, 0}, SquareIdentifier: A2, Piece: NewPawn(WhiteColor)},
		B2: Square{Empty: false, Coordinates: Coordinate{1, 1}, SquareIdentifier: B2, Piece: NewPawn(WhiteColor)},
		C2: Square{Empty: false, Coordinates: Coordinate{1, 2}, SquareIdentifier: C2, Piece: NewPawn(WhiteColor)},
		D2: Square{Empty: false, Coordinates: Coordinate{1, 3}, SquareIdentifier: D2, Piece: NewPawn(WhiteColor)},
		E2: Square{Empty: false, Coordinates: Coordinate{1, 4}, SquareIdentifier: E2, Piece: NewPawn(WhiteColor)},
		F2: Square{Empty: false, Coordinates: Coordinate{1, 5}, SquareIdentifier: F2, Piece: NewPawn(WhiteColor)},
		G2: Square{Empty: false, Coordinates: Coordinate{1, 6}, SquareIdentifier: G2, Piece: NewPawn(WhiteColor)},
		H2: Square{Empty: false, Coordinates: Coordinate{1, 7}, SquareIdentifier: H2, Piece: NewPawn(WhiteColor)},

		// Black pieces
		A8: Square{Empty: false, Coordinates: Coordinate{7, 0}, SquareIdentifier: A8, Piece: NewRook(BlackColor)},
		B8: Square{Empty: false, Coordinates: Coordinate{7, 1}, SquareIdentifier: B8, Piece: NewKnight(BlackColor)},
		C8: Square{Empty: false, Coordinates: Coordinate{7, 2}, SquareIdentifier: C8, Piece: NewBishop(BlackColor)},
		D8: Square{Empty: false, Coordinates: Coordinate{7, 3}, SquareIdentifier: D8, Piece: NewQueen(BlackColor)},
		E8: Square{Empty: false, Coordinates: Coordinate{7, 4}, SquareIdentifier: E8, Piece: NewKing(BlackColor)},
		F8: Square{Empty: false, Coordinates: Coordinate{7, 5}, SquareIdentifier: F8, Piece: NewBishop(BlackColor)},
		G8: Square{Empty: false, Coordinates: Coordinate{7, 6}, SquareIdentifier: G8, Piece: NewKnight(BlackColor)},
		H8: Square{Empty: false, Coordinates: Coordinate{7, 7}, SquareIdentifier: H8, Piece: NewRook(BlackColor)},
		A7: Square{Empty: false, Coordinates: Coordinate{6, 0}, SquareIdentifier: A7, Piece: NewPawn(BlackColor)},
		B7: Square{Empty: false, Coordinates: Coordinate{6, 1}, SquareIdentifier: B7, Piece: NewPawn(BlackColor)},
		C7: Square{Empty: false, Coordinates: Coordinate{6, 2}, SquareIdentifier: C7, Piece: NewPawn(BlackColor)},
		D7: Square{Empty: false, Coordinates: Coordinate{6, 3}, SquareIdentifier: D7, Piece: NewPawn(BlackColor)},
		E7: Square{Empty: false, Coordinates: Coordinate{6, 4}, SquareIdentifier: E7, Piece: NewPawn(BlackColor)},
		F7: Square{Empty: false, Coordinates: Coordinate{6, 5}, SquareIdentifier: F7, Piece: NewPawn(BlackColor)},
		G7: Square{Empty: false, Coordinates: Coordinate{6, 6}, SquareIdentifier: G7, Piece: NewPawn(BlackColor)},
		H7: Square{Empty: false, Coordinates: Coordinate{6, 7}, SquareIdentifier: H7, Piece: NewPawn(BlackColor)},

		// Empty pieces
		A3: Square{Empty: true, Coordinates: Coordinate{2, 0}, SquareIdentifier: A3},
		B3: Square{Empty: true, Coordinates: Coordinate{2, 1}, SquareIdentifier: B3},
		C3: Square{Empty: true, Coordinates: Coordinate{2, 2}, SquareIdentifier: C3},
		D3: Square{Empty: true, Coordinates: Coordinate{2, 3}, SquareIdentifier: D3},
		E3: Square{Empty: true, Coordinates: Coordinate{2, 4}, SquareIdentifier: E3},
		F3: Square{Empty: true, Coordinates: Coordinate{2, 5}, SquareIdentifier: F3},
		G3: Square{Empty: true, Coordinates: Coordinate{2, 6}, SquareIdentifier: G3},
		H3: Square{Empty: true, Coordinates: Coordinate{2, 7}, SquareIdentifier: H3},

		A4: Square{Empty: true, Coordinates: Coordinate{3, 0}, SquareIdentifier: A4},
		B4: Square{Empty: true, Coordinates: Coordinate{3, 1}, SquareIdentifier: B4},
		C4: Square{Empty: true, Coordinates: Coordinate{3, 2}, SquareIdentifier: C4},
		D4: Square{Empty: true, Coordinates: Coordinate{3, 3}, SquareIdentifier: D4},
		E4: Square{Empty: true, Coordinates: Coordinate{3, 4}, SquareIdentifier: E4},
		F4: Square{Empty: true, Coordinates: Coordinate{3, 5}, SquareIdentifier: F4},
		G4: Square{Empty: true, Coordinates: Coordinate{3, 6}, SquareIdentifier: G4},
		H4: Square{Empty: true, Coordinates: Coordinate{3, 7}, SquareIdentifier: H4},

		A5: Square{Empty: true, Coordinates: Coordinate{4, 0}, SquareIdentifier: A5},
		B5: Square{Empty: true, Coordinates: Coordinate{4, 1}, SquareIdentifier: B5},
		C5: Square{Empty: true, Coordinates: Coordinate{4, 2}, SquareIdentifier: C5},
		D5: Square{Empty: true, Coordinates: Coordinate{4, 3}, SquareIdentifier: D5},
		E5: Square{Empty: true, Coordinates: Coordinate{4, 4}, SquareIdentifier: E5},
		F5: Square{Empty: true, Coordinates: Coordinate{4, 5}, SquareIdentifier: F5},
		G5: Square{Empty: true, Coordinates: Coordinate{4, 6}, SquareIdentifier: G5},
		H5: Square{Empty: true, Coordinates: Coordinate{4, 7}, SquareIdentifier: H5},

		A6: Square{Empty: true, Coordinates: Coordinate{5, 0}, SquareIdentifier: A6},
		B6: Square{Empty: true, Coordinates: Coordinate{5, 1}, SquareIdentifier: B6},
		C6: Square{Empty: true, Coordinates: Coordinate{5, 2}, SquareIdentifier: C6},
		D6: Square{Empty: true, Coordinates: Coordinate{5, 3}, SquareIdentifier: D6},
		E6: Square{Empty: true, Coordinates: Coordinate{5, 4}, SquareIdentifier: E6},
		F6: Square{Empty: true, Coordinates: Coordinate{5, 5}, SquareIdentifier: F6},
		G6: Square{Empty: true, Coordinates: Coordinate{5, 6}, SquareIdentifier: G6},
		H6: Square{Empty: true, Coordinates: Coordinate{5, 7}, SquareIdentifier: H6},
	}
}

// EatPiece removes a piece from Board location, subtract it from player
func (b *board) EatPiece(loc SquareIdentifier) Piece {
	// TODO: WIP
	// SHould remove piece game logic be located here ?

	var piece Piece
	square := b.squares[loc]
	square.Empty = true
	piece = square.Piece
	// 	square.Piece = Piece{}
	return piece
}

func (b *board) Squares() Squares {
	return b.squares
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

// Coordinate a coordinate within the Board, starting at 0,0 = A1, 8,8 = H8
type Coordinate struct {
	X uint8
	Y uint8
}
