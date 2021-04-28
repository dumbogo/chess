package engine

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewKing(t *testing.T) {
	assert := assert.New(t)
	king := NewKing(WhiteColor)
	assert.Equal(KingIdentifier, king.Identifier())
	assert.Equal(WhiteColor, king.Color())
}

func TestKingCanMove(t *testing.T) {
	assert := assert.New(t)

	whiteKing := NewKing(WhiteColor)
	movements := []Movement{}

	ctrlMockBoard := gomock.NewController(t)
	defer ctrlMockBoard.Finish()
	mockBoard := NewMockBoard(ctrlMockBoard)

	a1 := Square{
		Empty:            true,
		Coordinates:      Coordinate{0, 0},
		SquareIdentifier: A1,
	}
	b2 := Square{
		Empty:            false,
		Piece:            whiteKing,
		Coordinates:      Coordinate{1, 1},
		SquareIdentifier: B2,
	}
	c3 := Square{
		Empty:            true,
		Coordinates:      Coordinate{2, 2},
		SquareIdentifier: C3,
	}

	ctrlBlackPieceD4 := gomock.NewController(t)
	defer ctrlBlackPieceD4.Finish()
	mockBlackPieceD4 := NewMockPiece(ctrlBlackPieceD4)
	d4 := Square{
		Empty:            false,
		Piece:            mockBlackPieceD4,
		Coordinates:      Coordinate{3, 3},
		SquareIdentifier: D4,
	}

	squares := map[SquareIdentifier]Square{
		A1: a1, // Empty
		B2: b2, // WhiteKing
		C3: c3, // Empty
		D4: d4, // BlackPiece
	}
	mockBoard.
		EXPECT().
		Squares().
		Return(squares)

	mockBlackPieceD4.
		EXPECT().
		Color().
		Return(BlackColor)
	mockBlackPieceD4.
		EXPECT().
		CanMove(mockBoard, movements, d4, c3).
		Return(true)

	// Cannot move due to king being eaten at c3 by mockBlackPieceD4
	assert.False(whiteKing.CanMove(mockBoard, movements, b2, c3))

	mockBoard.
		EXPECT().
		Squares().
		Return(squares)
	mockBlackPieceD4.
		EXPECT().
		Color().
		Return(BlackColor)
	mockBlackPieceD4.
		EXPECT().
		CanMove(mockBoard, movements, d4, a1).
		Return(false)

	// Valid movement
	assert.True(whiteKing.CanMove(mockBoard, movements, b2, a1))
}

func TestNewQueen(t *testing.T) {
	assert := assert.New(t)
	p := NewQueen(BlackColor)
	assert.Equal(QueenIdentifier, p.Identifier())
	assert.Equal(BlackColor, p.Color())
}

func TestQueenCanMove(t *testing.T) { // TODO: End testcase TestQueenCanMove
	t.Skip()
}

func TestNewRook(t *testing.T) {
	assert := assert.New(t)
	p := NewRook(BlackColor)
	assert.Equal(RookIdentifier, p.Identifier())
	assert.Equal(BlackColor, p.Color())
}

func TestRookCanMove(t *testing.T) {
	assert := assert.New(t)
	whiteRook := NewRook(WhiteColor)
	movements := []Movement{}

	ctrlMockBoard := gomock.NewController(t)
	defer ctrlMockBoard.Finish()
	mockBoard := NewMockBoard(ctrlMockBoard)

	d4 := Square{
		Empty:            true,
		Coordinates:      Coordinate{3, 3},
		SquareIdentifier: D4,
	}
	d5 := Square{
		Empty:            true,
		Coordinates:      Coordinate{3, 4},
		SquareIdentifier: D5,
	}
	d6 := Square{
		Empty:            false,
		Coordinates:      Coordinate{3, 5},
		SquareIdentifier: D6,
	}
	d7 := Square{
		Empty:            false,
		Piece:            whiteRook,
		Coordinates:      Coordinate{3, 6},
		SquareIdentifier: D7,
	}

	squares := map[SquareIdentifier]Square{
		D4: d4,
		D5: d5,
		D6: d6,
		D7: d7,
	}
	mockBoard.
		EXPECT().
		Squares().
		Return(squares)
	assert.False(whiteRook.CanMove(mockBoard, movements, d7, d4))

	// Case: When "to" square has an eating piece
	ctrlMockBoard = gomock.NewController(t)
	defer ctrlMockBoard.Finish()
	mockBoard = NewMockBoard(ctrlMockBoard)

	ctrlMockPiece := gomock.NewController(t)
	blackPiece := NewMockPiece(ctrlMockPiece)

	d4 = Square{
		Empty:            false,
		Piece:            whiteRook,
		Coordinates:      Coordinate{3, 3},
		SquareIdentifier: D4,
	}
	d5 = Square{
		Empty:            true,
		Coordinates:      Coordinate{3, 4},
		SquareIdentifier: D5,
	}
	d6 = Square{
		Empty:            true,
		Coordinates:      Coordinate{3, 5},
		SquareIdentifier: D6,
	}
	d7 = Square{
		Empty:            false,
		Piece:            blackPiece,
		Coordinates:      Coordinate{3, 6},
		SquareIdentifier: D7,
	}
	squares = map[SquareIdentifier]Square{
		D4: d4,
		D5: d5,
		D6: d6,
		D7: d7,
	}
	mockBoard.
		EXPECT().
		Squares().
		Return(squares)

	blackPiece.
		EXPECT().
		Color().
		Return(BlackColor)
	assert.True(whiteRook.CanMove(mockBoard, movements, d4, d7))

	// Case: move on X axis,
	a1 := Square{
		Empty:            true,
		Coordinates:      Coordinate{0, 0},
		SquareIdentifier: A1,
	}
	b1 := Square{
		Empty:            true,
		Coordinates:      Coordinate{1, 0},
		SquareIdentifier: B1,
	}
	c1 := Square{
		Empty:            true,
		Coordinates:      Coordinate{1, 0},
		SquareIdentifier: C1,
	}
	d1 := Square{
		Empty:            false,
		Piece:            whiteRook,
		Coordinates:      Coordinate{1, 0},
		SquareIdentifier: C1,
	}

	squares = map[SquareIdentifier]Square{
		A1: a1,
		B1: b1,
		C1: c1,
		D1: d1,
	}
	mockBoard.
		EXPECT().
		Squares().
		Return(squares)
	assert.True(whiteRook.CanMove(mockBoard, movements, d1, a1))
}

func TestNewBishop(t *testing.T) {
	assert := assert.New(t)
	p := NewBishop(BlackColor)
	assert.Equal(BishopIdentifier, p.Identifier())
	assert.Equal(BlackColor, p.Color())
}

func TestBishopCanMove(t *testing.T) {
	assert := assert.New(t)
	blackBishop := NewBishop(BlackColor)
	movements := []Movement{}

	// Case: when no pieces between, can do movement
	ctrlMockBoard := gomock.NewController(t)
	defer ctrlMockBoard.Finish()
	mockBoard := NewMockBoard(ctrlMockBoard)

	a1 := Square{Empty: true, Coordinates: Coordinate{0, 0}, SquareIdentifier: A1}
	b2 := Square{Empty: true, Coordinates: Coordinate{1, 1}, SquareIdentifier: B2}
	c3 := Square{Empty: false, Coordinates: Coordinate{2, 2}, SquareIdentifier: C3, Piece: blackBishop}

	squares := Squares{
		A1: a1,
		B2: b2,
		C3: c3,
	}

	mockBoard.
		EXPECT().
		Squares().
		Return(squares)
	assert.True(blackBishop.CanMove(mockBoard, movements, c3, a1))

	// Case: when there is a piece between from, to, cannot move
	ctrlMockBoard = gomock.NewController(t)
	defer ctrlMockBoard.Finish()
	mockBoard = NewMockBoard(ctrlMockBoard)

	a1 = Square{Empty: false, Coordinates: Coordinate{0, 0}, SquareIdentifier: A1, Piece: blackBishop}
	// Only Empty flag matters, set to false
	b2 = Square{Empty: false, Coordinates: Coordinate{1, 1}, SquareIdentifier: B2}
	c3 = Square{Empty: true, Coordinates: Coordinate{2, 2}, SquareIdentifier: C3}

	squares = Squares{
		A1: a1,
		B2: b2,
		C3: c3,
	}

	mockBoard.
		EXPECT().
		Squares().
		Return(squares)
	assert.False(blackBishop.CanMove(mockBoard, movements, a1, c3))
}

func TestNewKnight(t *testing.T) {
	assert := assert.New(t)
	p := NewKnight(BlackColor)
	assert.Equal(KnightIdentifier, p.Identifier())
	assert.Equal(BlackColor, p.Color())
}

func TestKnightCanMove(t *testing.T) {
	assert := assert.New(t)
	blackKnight := NewKnight(BlackColor)

	movements := []Movement{}

	ctrlMockBoard := gomock.NewController(t)
	defer ctrlMockBoard.Finish()
	mockBoard := NewMockBoard(ctrlMockBoard)

	c3 := Square{Empty: false, Coordinates: Coordinate{2, 2}, SquareIdentifier: C3, Piece: blackKnight}
	b5 := Square{Empty: true, Coordinates: Coordinate{1, 4}, SquareIdentifier: B5}
	d5 := Square{Empty: true, Coordinates: Coordinate{3, 4}, SquareIdentifier: D5}
	e4 := Square{Empty: true, Coordinates: Coordinate{4, 3}, SquareIdentifier: E4}
	e2 := Square{Empty: false, Coordinates: Coordinate{4, 1}, SquareIdentifier: E2, Piece: NewPawn(WhiteColor)}
	d1 := Square{Empty: false, Coordinates: Coordinate{3, 0}, SquareIdentifier: D1, Piece: NewPawn(BlackColor)}
	b1 := Square{Empty: true, Coordinates: Coordinate{1, 0}, SquareIdentifier: B1}
	a2 := Square{Empty: true, Coordinates: Coordinate{0, 1}, SquareIdentifier: A2}
	a4 := Square{Empty: true, Coordinates: Coordinate{0, 3}, SquareIdentifier: A4}

	assert.True(blackKnight.CanMove(mockBoard, movements, c3, b5))
	assert.True(blackKnight.CanMove(mockBoard, movements, c3, d5))
	assert.True(blackKnight.CanMove(mockBoard, movements, c3, e4))
	assert.True(blackKnight.CanMove(mockBoard, movements, c3, e4))
	assert.True(blackKnight.CanMove(mockBoard, movements, c3, e2))
	assert.False(blackKnight.CanMove(mockBoard, movements, c3, d1))
	assert.True(blackKnight.CanMove(mockBoard, movements, c3, b1))
	assert.True(blackKnight.CanMove(mockBoard, movements, c3, a2))
	assert.True(blackKnight.CanMove(mockBoard, movements, c3, a4))
	assert.False(blackKnight.CanMove(mockBoard, movements, a4, c3))
}

func TestNewPawn(t *testing.T) {
	assert := assert.New(t)
	pawn := NewPawn(BlackColor)
	assert.Equal(PawnIdentifier, pawn.Identifier())
	assert.Equal(BlackColor, pawn.Color())
}

func TestPawnCanMove(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// White movements validations
	whitePawn := NewPawn(WhiteColor)
	blackPawn := NewPawn(BlackColor)
	// when first movement, advance two spaces
	movements := []Movement{}
	square1 := Square{
		Empty:            false,
		Coordinates:      Coordinate{0, 1},
		SquareIdentifier: A2,
		Piece:            whitePawn,
	}
	square2 := Square{
		Empty:            true,
		Coordinates:      Coordinate{0, 3},
		SquareIdentifier: A4,
	}
	board := NewMockBoard(ctrl)
	assert.True(whitePawn.CanMove(board, movements, square1, square2))

	// when moving to empty
	movements = []Movement{}
	square1 = Square{
		Empty:            false,
		Piece:            whitePawn,
		Coordinates:      Coordinate{3, 3},
		SquareIdentifier: D4,
	}
	square2 = Square{
		Empty:            true,
		Coordinates:      Coordinate{3, 4},
		SquareIdentifier: D5,
	}
	board = NewMockBoard(ctrl)
	assert.True(whitePawn.CanMove(board, movements, square1, square2))

	// Eating a piece
	movements = []Movement{}
	square1 = Square{
		Empty:            false,
		Piece:            whitePawn,
		Coordinates:      Coordinate{3, 3},
		SquareIdentifier: D4,
	}
	square2 = Square{
		Empty:            false,
		Piece:            blackPawn,
		Coordinates:      Coordinate{4, 4},
		SquareIdentifier: E5,
	}
	board = NewMockBoard(ctrl)
	assert.True(whitePawn.CanMove(board, movements, square1, square2))

	// Black movements
	// when first movement, advance two spaces
	movements = []Movement{}
	square1 = Square{
		Empty:            false,
		Coordinates:      Coordinate{5, 6},
		SquareIdentifier: F7,
		Piece:            blackPawn,
	}
	square2 = Square{
		Empty:            true,
		Coordinates:      Coordinate{5, 4},
		SquareIdentifier: F5,
	}
	board = NewMockBoard(ctrl)
	assert.True(blackPawn.CanMove(board, movements, square1, square2))

	// when moving to empty
	movements = []Movement{}
	square1 = Square{
		Empty:            false,
		Piece:            blackPawn,
		Coordinates:      Coordinate{3, 3},
		SquareIdentifier: D4,
	}
	square2 = Square{
		Empty:            true,
		Coordinates:      Coordinate{3, 2},
		SquareIdentifier: D3,
	}
	board = NewMockBoard(ctrl)
	assert.True(whitePawn.CanMove(board, movements, square1, square2))

	// Eating a piece
	movements = []Movement{}
	square1 = Square{
		Empty:            false,
		Piece:            blackPawn,
		Coordinates:      Coordinate{3, 3},
		SquareIdentifier: D4,
	}
	square2 = Square{
		Empty:            false,
		Piece:            whitePawn,
		Coordinates:      Coordinate{2, 2},
		SquareIdentifier: C3,
	}
	board = NewMockBoard(ctrl)
	assert.True(whitePawn.CanMove(board, movements, square1, square2))
}
