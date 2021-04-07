package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBoard(t *testing.T) {
	assert := assert.New(t)
	p1, p2 := Player{Name: "Hugo"}, Player{Name: "Pedro"}
	board := NewBoard(&p1, &p2)
	assert.Equal("Hugo", board.WhitePlayer().Name)
	assert.Equal("Pedro", board.BlackPlayer().Name)
	// TODO: Add assertions
}

// TODO: testcases for:
// 	FillSquares()
// 	EatPiece(loc SquareIdentifier) Piece
// 	WhitePlayer() *Player
// 	BlackPlayer() *Player
// 	Squares() Squares
