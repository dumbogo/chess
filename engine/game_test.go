package engine

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var testPlayerBlack = Player{Name: "Joel", Color: BlackColor}
var testPlayerWhite = Player{Name: "Luis", Color: WhiteColor}

var testCaseGameGenerate = func() (Game, error) {
	return NewGame(
		"asaditosGame",
		testPlayerBlack,
		testPlayerWhite,
	)
}

func TestNewGame(t *testing.T) {
	// When error returns
	assert := assert.New(t)
	game, err := NewGame(
		"asaditosGame",
		Player{Name: "Joel", Color: WhiteColor},
		Player{Name: "Luis", Color: WhiteColor},
	)
	assert.Equal(nil, game)
	assert.Equal(errors.New("must define black and white players"), err)

	// When no errors
	_, err = testCaseGameGenerate()
	assert.Equal(nil, err)
}

func TestLoadGame(t *testing.T) {
	board := NewBoard(&testPlayerWhite, &testPlayerBlack)
	blackPieces := map[PieceIdentifier]uint8{}
	whitePieces := map[PieceIdentifier]uint8{}
	LoadGame(
		"wipgame",
		board,
		testPlayerWhite,
		testPlayerWhite,
		testPlayerBlack,
		whitePieces,
		blackPieces,
		make([]Movement, 0),
	)
}

func TestTurn(t *testing.T) {
	assert := assert.New(t)
	testCaseGame, _ := testCaseGameGenerate()
	turn := testCaseGame.Turn()
	assert.Equal(testPlayerWhite, turn)
}

func TestMove(t *testing.T) {
	var testCaseGame Game
	assert := assert.New(t)

	testCaseGame, _ = testCaseGameGenerate()
	ok, e := testCaseGame.Move(testPlayerWhite, A2, A3)
	assert.Equal(true, ok)
	assert.Empty(e)
	assert.Equal(testPlayerBlack, testCaseGame.Turn())
	// This is tooo deph, need to shorten method to return squares
	assert.Equal(true, testCaseGame.Board().Squares()[A2].Empty)
	assert.Empty(testCaseGame.Board().Squares()[A2].Piece)
	assert.Equal(false, testCaseGame.Board().Squares()[A3].Empty)
	assert.Equal(PawnIdentifier, testCaseGame.Board().Squares()[A3].Piece.Identifier())

	// Try to move an empty square
	ok, e = testCaseGame.Move(testPlayerBlack, A4, A5)
	assert.Equal(false, ok)
	assert.NotNil(e)
	// Move pawn
	ok, e = testCaseGame.Move(testPlayerBlack, A7, A5)
	assert.Equal(true, ok)
	assert.Empty(e)

	// whitep tries to move black pawn
	ok, e = testCaseGame.Move(testPlayerBlack, B7, B5)
	assert.Equal(false, ok)
	assert.NotNil(e)

	// Testcase when movement eats a piece, white pawn eats black pawn
	testCaseGame, _ = testCaseGameGenerate()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c3whitePawnSquare := Square{
		Empty:            false,
		Piece:            NewPawn(WhiteColor),
		SquareIdentifier: C3,
		Coordinates: Coordinate{
			X: 3,
			Y: 2,
		},
	}
	d4blackPawnSquare := Square{
		Empty:            false,
		Piece:            NewPawn(BlackColor),
		SquareIdentifier: D4,
		Coordinates: Coordinate{
			X: 4,
			Y: 3,
		},
	}

	testCaseGame.Board().Squares()[C3] = c3whitePawnSquare
	testCaseGame.Board().Squares()[D4] = d4blackPawnSquare
	ok, err := testCaseGame.Move(testPlayerWhite, C3, D4)
	assert.True(ok)
	assert.Empty(err)
}

func TestIsCheckmateBy(t *testing.T) {
	assert := assert.New(t)

	// Case: checkmate, Fool's mate
	Joel := Player{"Joel", WhiteColor}
	Luis := Player{"Luis", BlackColor}
	game, e := NewGame("Fool's mate", Luis, Joel)
	check(e)

	game.Move(Joel, E2, E4)
	game.Move(Luis, G7, G5)
	game.Move(Joel, B1, C3)

	no := game.IsCheckmateBy(Joel)
	assert.False(no)
	game.Move(Luis, F7, F5)
	game.Move(Joel, D1, H5)
	yes := game.IsCheckmateBy(Joel)
	assert.True(yes)

	ok, err := game.Move(Luis, C7, C6)
	assert.False(ok)
	assert.Error(err)
}

func TestIsCheck(t *testing.T) {
	assert := assert.New(t)

	Joel := Player{"Joel", WhiteColor}
	Luis := Player{"Luis", BlackColor}
	game, e := NewGame("rapid check", Luis, Joel)
	assert.Nil(e)

	game.Move(Joel, D2, D4)
	game.Move(Luis, E7, E5)
	game.Move(Joel, D1, D3)
	game.Move(Luis, A7, A6)
	game.Move(Joel, D3, E3)
	game.Move(Luis, A6, A5)
	game.Move(Joel, E3, E5)
	ok, _ := game.Move(Luis, A5, A4)
	assert.False(ok)
	game.Move(Luis, F8, E7)
	t.Skip()
	// assert.True(ok)
}

func TestRollback(t *testing.T) {
	assert := assert.New(t)
	Joel := Player{"Joel", WhiteColor}
	Luis := Player{"Luis", BlackColor}
	game, e := NewGame("rapid check", Luis, Joel)
	assert.Nil(e)

	game.Move(Joel, D2, D4)
	game.Move(Luis, E7, E5)
	game.Rollback(1)
	assert.Equal(1, len(game.Movements()))
	assert.EqualValues(
		Square{
			Empty:            true,
			Piece:            nil,
			Coordinates:      SquareIdentifierToCoordinate(E5),
			SquareIdentifier: E5,
		},
		game.Board().Squares()[E5],
	)

	assert.EqualValues(
		Square{
			Empty:            false,
			Piece:            NewPawn(BlackColor),
			Coordinates:      SquareIdentifierToCoordinate(E7),
			SquareIdentifier: E7,
		},
		game.Board().Squares()[E7],
	)
	assert.EqualValues(Luis, game.Turn())
}

func TestBoard(t *testing.T) {
	// TODO:WIP
	t.Skip()
}

func TestGameString(t *testing.T) {
	assert := assert.New(t)
	testCaseGame, _ := testCaseGameGenerate()
	r := testCaseGame.String()
	expected := `+---+----+----+----+----+----+----+----+----+
| 8 | BR | Bk | BB | BQ | BK | BB | Bk | BR |
| 7 | BP | BP | BP | BP | BP | BP | BP | BP |
| 6 |    |    |    |    |    |    |    |    |
| 5 |    |    |    |    |    |    |    |    |
| 4 |    |    |    |    |    |    |    |    |
| 3 |    |    |    |    |    |    |    |    |
| 2 | WP | WP | WP | WP | WP | WP | WP | WP |
| 1 | WR | Wk | WB | WQ | WK | WB | Wk | WR |
|   | a  | b  | c  | d  | e  | f  | g  | h  |
+---+----+----+----+----+----+----+----+----+
`
	assert.Equal(expected, r)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
