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
	t.Skip()
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pawn := NewPawn(BlackColor)

	// when first movement, advance two spaces
	movements := []Movement{}
	square1 := Square{
		Empty:            false,
		Coordinates:      Coordinate{5, 1, F2},
		SquareIdentifier: F2,
		Piece:            pawn,
	}
	square2 := Square{
		Empty:            true,
		Coordinates:      Coordinate{5, 3, F4},
		SquareIdentifier: F4,
	}
	board := NewMockBoard(ctrl)
	board.
		EXPECT().
		Squares().
		Return(Squares{F2: square1, F4: square2})
	assert.True(pawn.CanMove(board, movements, square1, square2))

	// when moving to empty
	movements = []Movement{}
	square1 = Square{
		Empty:            false,
		Piece:            pawn,
		Coordinates:      Coordinate{1, 3, D3},
		SquareIdentifier: D3,
	}
	square2 = Square{
		Empty:            true,
		Coordinates:      Coordinate{1, 4, B5},
		SquareIdentifier: B5,
	}
	board = NewMockBoard(ctrl)
	board.
		EXPECT().
		Squares().
		Return(Squares{D3: square1, B5: square2})

	ok := pawn.CanMove(board, movements, square1, square2)
	assert.True(ok)

	// TODO: WIP
	// En passant
	// Eating a piece
	// Illegal movement: blocked by other piece
}

// TODO: add tests CanMove interface method for each piece
