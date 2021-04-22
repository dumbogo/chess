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

// TODO: add tests CanMove interface method for each piece
