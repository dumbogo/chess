package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var p1, p2 Player = Player{Name: "Hugo"}, Player{Name: "Pedro"}
var testableBoard = func() Board {
	return NewBoard(&p1, &p2)
}

func TestNewBoard(t *testing.T) {
	assert := assert.New(t)
	board := testableBoard()
	assert.Equal("Hugo", board.WhitePlayer().Name)
	assert.Equal("Pedro", board.BlackPlayer().Name)

	squares := board.Squares()
	assert.Equal(false, squares[A1].Empty)
	assert.Equal(false, squares[B1].Empty)
	assert.Equal(false, squares[C1].Empty)
	assert.Equal(false, squares[D1].Empty)
	assert.Equal(false, squares[E1].Empty)
	assert.Equal(false, squares[F1].Empty)
	assert.Equal(false, squares[G1].Empty)
	assert.Equal(false, squares[H1].Empty)

	assert.Equal(false, squares[A2].Empty)
	assert.Equal(false, squares[B2].Empty)
	assert.Equal(false, squares[C2].Empty)
	assert.Equal(false, squares[D2].Empty)
	assert.Equal(false, squares[E2].Empty)
	assert.Equal(false, squares[F2].Empty)
	assert.Equal(false, squares[G2].Empty)
	assert.Equal(false, squares[H2].Empty)

	assert.Equal(false, squares[A7].Empty)
	assert.Equal(false, squares[B7].Empty)
	assert.Equal(false, squares[C7].Empty)
	assert.Equal(false, squares[D7].Empty)
	assert.Equal(false, squares[E7].Empty)
	assert.Equal(false, squares[F7].Empty)
	assert.Equal(false, squares[G7].Empty)
	assert.Equal(false, squares[H7].Empty)

	assert.Equal(false, squares[A8].Empty)
	assert.Equal(false, squares[B8].Empty)
	assert.Equal(false, squares[C8].Empty)
	assert.Equal(false, squares[D8].Empty)
	assert.Equal(false, squares[E8].Empty)
	assert.Equal(false, squares[F8].Empty)
	assert.Equal(false, squares[G8].Empty)
	assert.Equal(false, squares[H8].Empty)

	for s := A3; s < A7; s++ {
		assert.Equal(true, squares[s].Empty)
	}
}

func TestEatPiece(t *testing.T) {
	assert := assert.New(t)
	board := testableBoard()

	piece := board.EatPiece(H7)
	square := board.Squares()[H7]

	assert.Equal(NewPawn(BlackColor), piece)
	assert.Equal(square.Empty, true)
	assert.Empty(square.Piece)
}

func TestFillSquare(t *testing.T) {
	assert := assert.New(t)
	board := testableBoard()
	whitePawn := NewPawn(WhiteColor)
	board.FillSquare(D5, whitePawn)
	square := board.Squares()[D5]
	assert.EqualValues(
		Square{
			Empty:            false,
			Piece:            whitePawn,
			Coordinates:      Coordinate{3, 4},
			SquareIdentifier: D5,
		},
		square,
	)
}

func TestCoordinateToSquareIdentifier(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(A1, CoordinateToSquareIdentifier(Coordinate{0, 0}))
	assert.Equal(H2, CoordinateToSquareIdentifier(Coordinate{7, 1}))
	assert.Equal(H6, CoordinateToSquareIdentifier(Coordinate{7, 5}))
	assert.Equal(H8, CoordinateToSquareIdentifier(Coordinate{7, 7}))
}

func TestSquareIdentifierToCoordinate(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(Coordinate{0, 0}, SquareIdentifierToCoordinate(A1))
	assert.Equal(Coordinate{7, 1}, SquareIdentifierToCoordinate(H2))
	assert.Equal(Coordinate{7, 5}, SquareIdentifierToCoordinate(H6))
	assert.Equal(Coordinate{7, 7}, SquareIdentifierToCoordinate(H8))
}
