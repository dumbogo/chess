package engine

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testPlayerBlack = Player{Name: "Joel", Color: BlackColor}
var testPlayerWhite = Player{Name: "Luis", Color: WhiteColor}

var testCaseGame, eGame = NewGame(
	"asaditosGame",
	testPlayerBlack,
	testPlayerWhite,
)

func TestNewGame(t *testing.T) {
	// When error returns
	assert := assert.New(t)
	game, err := NewGame(
		"asaditosGame",
		Player{Name: "Joel", Color: WhiteColor},
		Player{Name: "Luis", Color: WhiteColor},
	)
	assert.Equal(nil, game)
	assert.Equal(errors.New("Must define black and white players"), err)
	// When no errors
	assert.Equal(nil, eGame)
}

func TestTurn(t *testing.T) {
	assert := assert.New(t)
	turn := testCaseGame.Turn()
	assert.Equal(testPlayerWhite, turn)
}

func TestMove(t *testing.T) {
	assert := assert.New(t)
	ok, e := testCaseGame.Move(testPlayerWhite, A2, A3)
	assert.Equal(true, ok)
	assert.Empty(e)
	assert.Equal(testPlayerBlack, testCaseGame.Turn())
	// This is tooo deph, need to shorten method to return squares
	assert.Equal(true, testCaseGame.Board().Squares()[A2].Empty)
	assert.Equal(nil, testCaseGame.Board().Squares()[A2].Piece)
	assert.Equal(false, testCaseGame.Board().Squares()[A3].Empty)
	assert.Equal(PawnIdentifier, testCaseGame.Board().Squares()[A3].Piece.Identifier())
}

func TestIsCheckBy(t *testing.T) {
	// TODO: wip
	t.Skip()
}

func TestIsCheckmate(t *testing.T) {
	// TODO: wip
	t.Skip()
}

func TestBoard(t *testing.T) {
	// TODO:WIP
	t.Skip()
}
