package engine

import (
	"errors"
)

// Player gammer
type Player struct {
	Name  string
	Color Color
}

// Game playable game
type Game interface {
	// Turn returns player turn
	Turn() Player
	// Move moves a piece in the Board, returns true if moved
	Move(player Player, form, to SquareIdentifier) (bool, error)
	// IsCheckBy returns Player making check
	IsCheckBy() Player
	// IsCheckmateBy returns Player making checkmate
	IsCheckmateBy() Player
	// Board get board
	Board() Board
}

// NewGame creates new Game
func NewGame(name string, black, white Player) (Game, error) {
	if white.Color != WhiteColor || black.Color != BlackColor {
		return nil, errors.New("Must define black and white players")
	}
	return &game{
		name:  name,
		board: NewBoard(&white, &black),
		turn:  white,
		white: white,
		black: black,
	}, nil
}

type game struct {
	name  string
	board Board
	turn  Player
	white Player
	black Player
}

func (g *game) Turn() Player {
	return g.turn
}

func (g *game) Move(player Player, from, to SquareIdentifier) (bool, error) {
	squareFrom := g.board.Squares()[from]
	squareTo := g.board.Squares()[to]
	if squareFrom.Empty == true {
		return false, nil
	}

	pieceToMove := squareFrom.Piece
	if pieceToMove.Color() != g.Turn().Color {
		return false, nil
	}
	canMove := squareFrom.Piece.CanMove(squareFrom, squareTo)
	if canMove == false {
		return false, nil
	}

	if !squareTo.Empty {
		pieceToEat := squareTo.Piece
		if pieceToEat.Color() == g.Turn().Color {
			return false, nil
		}
		g.board.EatPiece(to)
	}

	squareTo.Piece = pieceToMove
	squareTo.Empty = false
	squareFrom.Piece = nil
	squareFrom.Empty = true
	g.board.Squares()[to] = squareTo
	g.board.Squares()[from] = squareFrom
	g.changeTurn()
	return true, nil
}

func (g *game) IsCheckBy() Player {
	// TODO: End
	return g.white
}

func (g *game) IsCheckmateBy() Player {
	// TODO: end
	return g.black
}

func (g *game) Board() Board {
	return g.board
}

// String returns ASCII representation of the game
func (g *game) String() string {
	// TODO: End
	return ""
}

func (g *game) changeTurn() {
	if g.Turn().Color == WhiteColor {
		g.turn = g.black
	} else {
		g.turn = g.white
	}
}
