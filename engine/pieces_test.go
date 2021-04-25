package engine

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

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
		Coordinates:      Coordinate{0, 1, A2},
		SquareIdentifier: A2,
		Piece:            whitePawn,
	}
	square2 := Square{
		Empty:            true,
		Coordinates:      Coordinate{0, 3, A4},
		SquareIdentifier: A4,
	}
	board := NewMockBoard(ctrl)
	assert.True(whitePawn.CanMove(board, movements, square1, square2))

	// when moving to empty
	movements = []Movement{}
	square1 = Square{
		Empty:            false,
		Piece:            whitePawn,
		Coordinates:      Coordinate{3, 3, D4},
		SquareIdentifier: D4,
	}
	square2 = Square{
		Empty:            true,
		Coordinates:      Coordinate{3, 4, D5},
		SquareIdentifier: D5,
	}
	board = NewMockBoard(ctrl)
	assert.True(whitePawn.CanMove(board, movements, square1, square2))

	// Eating a piece
	movements = []Movement{}
	square1 = Square{
		Empty:            false,
		Piece:            whitePawn,
		Coordinates:      Coordinate{3, 3, D4},
		SquareIdentifier: D4,
	}
	square2 = Square{
		Empty:            false,
		Piece:            blackPawn,
		Coordinates:      Coordinate{4, 4, E5},
		SquareIdentifier: E5,
	}
	board = NewMockBoard(ctrl)
	assert.True(whitePawn.CanMove(board, movements, square1, square2))

	// Black movements
	// when first movement, advance two spaces
	movements = []Movement{}
	square1 = Square{
		Empty:            false,
		Coordinates:      Coordinate{5, 6, F7},
		SquareIdentifier: F7,
		Piece:            blackPawn,
	}
	square2 = Square{
		Empty:            true,
		Coordinates:      Coordinate{5, 4, F5},
		SquareIdentifier: F5,
	}
	board = NewMockBoard(ctrl)
	assert.True(blackPawn.CanMove(board, movements, square1, square2))

	// when moving to empty
	movements = []Movement{}
	square1 = Square{
		Empty:            false,
		Piece:            blackPawn,
		Coordinates:      Coordinate{3, 3, D4},
		SquareIdentifier: D4,
	}
	square2 = Square{
		Empty:            true,
		Coordinates:      Coordinate{3, 2, D3},
		SquareIdentifier: D3,
	}
	board = NewMockBoard(ctrl)
	assert.True(whitePawn.CanMove(board, movements, square1, square2))

	// Eating a piece
	movements = []Movement{}
	square1 = Square{
		Empty:            false,
		Piece:            blackPawn,
		Coordinates:      Coordinate{3, 3, D4},
		SquareIdentifier: D4,
	}
	square2 = Square{
		Empty:            false,
		Piece:            whitePawn,
		Coordinates:      Coordinate{2, 2, C3},
		SquareIdentifier: C3,
	}
	board = NewMockBoard(ctrl)
	assert.True(whitePawn.CanMove(board, movements, square1, square2))
}

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
		Coordinates:      Coordinate{0, 0, A1},
		SquareIdentifier: A1,
	}
	b2 := Square{
		Empty:            false,
		Piece:            whiteKing,
		Coordinates:      Coordinate{1, 1, B2},
		SquareIdentifier: B2,
	}
	c3 := Square{
		Empty:            true,
		Coordinates:      Coordinate{2, 2, C3},
		SquareIdentifier: C3,
	}

	ctrlBlackPieceD4 := gomock.NewController(t)
	defer ctrlBlackPieceD4.Finish()
	mockBlackPieceD4 := NewMockPiece(ctrlBlackPieceD4)
	d4 := Square{
		Empty:            false,
		Piece:            mockBlackPieceD4,
		Coordinates:      Coordinate{3, 3, D4},
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
		Coordinates:      Coordinate{3, 3, D4},
		SquareIdentifier: D4,
	}
	d5 := Square{
		Empty:            true,
		Coordinates:      Coordinate{3, 4, D5},
		SquareIdentifier: D5,
	}
	d6 := Square{
		Empty:            false,
		Coordinates:      Coordinate{3, 5, D6},
		SquareIdentifier: D6,
	}
	d7 := Square{
		Empty:            false,
		Piece:            whiteRook,
		Coordinates:      Coordinate{3, 6, D7},
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
		Coordinates:      Coordinate{3, 3, D4},
		SquareIdentifier: D4,
	}
	d5 = Square{
		Empty:            true,
		Coordinates:      Coordinate{3, 4, D5},
		SquareIdentifier: D5,
	}
	d6 = Square{
		Empty:            true,
		Coordinates:      Coordinate{3, 5, D6},
		SquareIdentifier: D6,
	}
	d7 = Square{
		Empty:            false,
		Piece:            blackPiece,
		Coordinates:      Coordinate{3, 6, D7},
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
		Coordinates:      Coordinate{0, 0, A1},
		SquareIdentifier: A1,
	}
	b1 := Square{
		Empty:            true,
		Coordinates:      Coordinate{1, 0, B1},
		SquareIdentifier: B1,
	}
	c1 := Square{
		Empty:            true,
		Coordinates:      Coordinate{1, 0, C1},
		SquareIdentifier: C1,
	}
	d1 := Square{
		Empty:            false,
		Piece:            whiteRook,
		Coordinates:      Coordinate{1, 0, C1},
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

// TODO: add tests CanMove interface method for each piece
