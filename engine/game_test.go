package engine

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
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
	assert.Equal(errors.New("Must define black and white players"), err)

	// When no errors
	game, err = testCaseGameGenerate()
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
	assert := assert.New(t)
	testCaseGame, _ := testCaseGameGenerate()
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
	assert.Empty(e)
	// Move pawn
	ok, e = testCaseGame.Move(testPlayerBlack, A7, A5)
	assert.Equal(true, ok)
	assert.Empty(e)

	// whitep tries to move black pawn
	ok, e = testCaseGame.Move(testPlayerBlack, B7, B5)
	assert.Equal(false, ok)
	assert.Empty(e)

	// Testcase when movement eats a piece(white eats black)
	testCaseGame, _ = testCaseGameGenerate()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	p1 := NewMockPiece(ctrl)
	p2 := NewMockPiece(ctrl)
	p1.
		EXPECT().
		Color().
		Return(WhiteColor)
	p2.
		EXPECT().
		Color().
		Return(BlackColor)
	p2.
		EXPECT().
		Identifier().
		Return(PawnIdentifier)

	square1 := Square{
		Empty:            false,
		Piece:            p1,
		SquareIdentifier: H4,
	}
	square2 := Square{
		Empty:            false,
		Piece:            p2,
		SquareIdentifier: H5,
	}
	p1.
		EXPECT().
		CanMove(testCaseGame.Board(), testCaseGame.Movements(), square1, square2).
		Return(true)

	testCaseGame.Board().Squares()[H4] = square1
	testCaseGame.Board().Squares()[H5] = square2
	ok, err := testCaseGame.Move(testPlayerBlack, H4, H5)
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
