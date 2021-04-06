package engine

// Board gameboard
type Board struct {
	WhitePlayer Player
	WhitePieces [16]Piece

	BlackPlayer Player
	BlackPieces [16]Piece
	Squares     Squares
}

// NewBoard creates a new Board to play
func NewBoard(whitePlayer, blackPlayer Player) Board {
	board := Board{
		BlackPlayer: blackPlayer,
		WhitePlayer: whitePlayer,
	}
	board.FillSquares()
	return board
}

// FillSquares fill squares board with new game
func (b *Board) FillSquares() {
	b.Squares = Squares{
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

		// TODO: end this
		A4: Square{Empty: true},
		B4: Square{Empty: true},
		C4: Square{Empty: true},
		D4: Square{Empty: true},
		E4: Square{Empty: true},
		F4: Square{Empty: true},
		G4: Square{Empty: true},
		H4: Square{Empty: true},

		A5: Square{Empty: true},
		B5: Square{Empty: true},
		C5: Square{Empty: true},
		D5: Square{Empty: true},
		E5: Square{Empty: true},
		F5: Square{Empty: true},
		G5: Square{Empty: true},
		H5: Square{Empty: true},

		A6: Square{Empty: true},
		B6: Square{Empty: true},
		C6: Square{Empty: true},
		D6: Square{Empty: true},
		E6: Square{Empty: true},
		F6: Square{Empty: true},
		G6: Square{Empty: true},
		H6: Square{Empty: true},
		// TODO: up to this point
	}
}

// Squares board squares, total of 8*8
type Squares map[SquareIdentifier]Square
