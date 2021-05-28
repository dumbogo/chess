package engine

import (
	"errors"
)

// TODO: Pawn Promotion
// TODO: check if En passant to delete Pawn
// TODO: check Castling(short & large)

// Player gammer
type Player struct {
	Name  string
	Color Color
}

// Movement a movement made
type Movement struct {
	Player     Player
	PieceMoved Piece
	PieceEaten Piece
	From       SquareIdentifier
	To         SquareIdentifier
}

// PiecesList is a map of pieces, each value contains number of pieces left in board
type PiecesList map[PieceIdentifier]uint8

// Game playable game
type Game interface {
	// Move moves a piece in the Board, returns true if moved
	Move(player Player, from, to SquareIdentifier) (bool, error)
	// Turn returns player turn
	Turn() Player
	// IsCheckBy returns true if Player makes check
	IsCheckBy(Player) bool
	// IsCheckmateBy returns true if Player makes checkmate
	IsCheckmateBy(Player) bool
	// Board get board
	Board() Board
	// Movements get all historic movements
	Movements() []Movement
	// Rollback returns to a previous stage, weight means how many steps back
	Rollback(weight int)

	WhitePieces() PiecesList
	BlackPieces() PiecesList

	String() string
}

type game struct {
	name  string
	board Board
	turn  Player
	white Player
	black Player

	whitePieces PiecesList
	blackPieces PiecesList
	movements   []Movement
}

// NewGame creates new Game
func NewGame(name string, black, white Player) (Game, error) {
	if white.Color != WhiteColor || black.Color != BlackColor {
		return nil, errors.New("must define black and white players")
	}

	blackPieces := PiecesList{
		RookIdentifier:   2,
		KnightIdentifier: 2,
		BishopIdentifier: 2,
		QueenIdentifier:  1,
		KingIdentifier:   1,
		PawnIdentifier:   8,
	}
	whitePieces := PiecesList{
		RookIdentifier:   2,
		KnightIdentifier: 2,
		BishopIdentifier: 2,
		QueenIdentifier:  1,
		KingIdentifier:   1,
		PawnIdentifier:   8,
	}
	return &game{
		name:        name,
		board:       NewBoard(&white, &black),
		turn:        white,
		white:       white,
		black:       black,
		blackPieces: blackPieces,
		whitePieces: whitePieces,
	}, nil
}

// LoadGame loads in game being played
func LoadGame(
	name string,
	board Board,
	turn Player,
	white, black Player,
	whitePieces PiecesList,
	blackPieces PiecesList,
	movements []Movement,
) (Game, error) {
	return &game{
		name:        name,
		board:       board,
		turn:        turn,
		white:       white,
		black:       black,
		whitePieces: whitePieces,
		blackPieces: blackPieces,
		movements:   movements,
	}, nil
}

func (g *game) Move(player Player, from, to SquareIdentifier) (bool, error) {
	if g.IsCheckmateBy(g.white) {
		return false, errors.New("checkmate, winner is white")
	}
	if g.IsCheckmateBy(g.black) {
		return false, errors.New("checkmate, winner is black")
	}
	squareFrom := g.board.Squares()[from]
	squareTo := g.board.Squares()[to]
	if squareFrom.Empty {
		return false, errors.New("square from is empty")
	}

	pieceToMove := squareFrom.Piece
	if pieceToMove.Color() != g.Turn().Color {
		return false, errors.New("piece is not for the player color")
	}

	canMove := squareFrom.Piece.CanMove(g.board, g.Movements(), squareFrom, squareTo)
	if !canMove {
		return false, errors.New("not valid piece movement")
	}

	var pieceEaten Piece
	if !squareTo.Empty {
		pieceEaten = g.board.EatPiece(to)
		g.removePiecePlayer(pieceEaten)
	}

	squareTo.Piece = pieceToMove
	squareTo.Empty = false
	squareFrom.Piece = nil
	squareFrom.Empty = true
	g.board.Squares()[to] = squareTo
	g.board.Squares()[from] = squareFrom
	g.movements = append(g.movements, Movement{
		Player:     player,
		PieceMoved: pieceToMove,
		PieceEaten: pieceEaten,
		From:       from,
		To:         to,
	})
	g.changeTurn()
	if g.IsCheckBy(g.Turn()) {
		g.Rollback(1)
		return false, errors.New("check")
	}
	return true, nil
}

func (g *game) Turn() Player {
	return g.turn
}

func (g *game) IsCheckBy(player Player) bool {
	var kingSquare Square
	playerColor := player.Color
	if playerColor == WhiteColor {
		kingSquare = getKingSquare(g.Board(), BlackColor)
	} else {
		kingSquare = getKingSquare(g.Board(), WhiteColor)
	}
	if kingEatableInSquare(kingSquare.Piece.Color(), g.Board(), g.Movements(), kingSquare) {
		return true
	}
	return false
}

func (g *game) IsCheckmateBy(player Player) bool {
	if !g.IsCheckBy(player) {
		return false
	}
	var kingSquare Square
	playerColor := player.Color
	if playerColor == WhiteColor {
		kingSquare = getKingSquare(g.Board(), BlackColor)
	} else {
		kingSquare = getKingSquare(g.Board(), WhiteColor)
	}
	// iterate all possible movements and checks if king can move to at least one direction
	for x := kingSquare.Coordinates.X - 1; x <= kingSquare.Coordinates.X+1; x++ {
		for y := kingSquare.Coordinates.Y - 1; y <= kingSquare.Coordinates.Y+1; y++ {
			if (x > MAXX || y > MAXY) || x == kingSquare.Coordinates.X && y == kingSquare.Coordinates.Y {
				continue
			}
			if kingSquare.Piece.CanMove(g.Board(), g.Movements(), kingSquare, g.Board().Squares()[CoordinateToSquareIdentifier(Coordinate{x, y})]) {
				return false
			}
		}
	}
	// TODO: review if any piece player can block the check
	// we probably want to have all pieces making check
	return true
}

func (g *game) Board() Board {
	return g.board
}

func (g *game) Movements() []Movement {
	return g.movements
}

func (g *game) Rollback(w int) {
	for i := 1; i <= w; i++ {
		lastMovement := g.movements[len(g.movements)-1]
		squareFrom := Square{
			Empty:            false,
			Piece:            lastMovement.PieceMoved,
			Coordinates:      SquareIdentifierToCoordinate(lastMovement.From),
			SquareIdentifier: lastMovement.From,
		}

		squareTo := Square{
			Piece:            lastMovement.PieceEaten,
			Coordinates:      SquareIdentifierToCoordinate(lastMovement.To),
			SquareIdentifier: lastMovement.To,
		}

		if squareTo.Piece == nil {
			squareTo.Empty = true
		} else {
			g.addPiecePlayer(squareTo.Piece)
		}
		g.board.Squares()[lastMovement.To] = squareTo
		g.board.Squares()[lastMovement.From] = squareFrom
		g.movements = g.movements[:len(g.movements)-1]
		g.changeTurn()
	}
}

func (g *game) WhitePieces() PiecesList {
	return g.whitePieces
}

func (g *game) BlackPieces() PiecesList {
	return g.blackPieces
}

// String returns ASCII representation of the game
func (g *game) String() string {
	return g.board.String()
}

func (g *game) changeTurn() {
	if g.Turn().Color == WhiteColor {
		g.turn = g.black
	} else {
		g.turn = g.white
	}
}

func (g *game) removePiecePlayer(p Piece) {
	color := p.Color()
	switch color {
	case BlackColor:
		g.blackPieces[p.Identifier()]--
	case WhiteColor:
		g.whitePieces[p.Identifier()]--
	}
}

func (g *game) addPiecePlayer(p Piece) {
	color := p.Color()
	switch color {
	case BlackColor:
		g.blackPieces[p.Identifier()]++
	case WhiteColor:
		g.whitePieces[p.Identifier()]++
	}
}

func getKingSquare(board Board, color Color) Square {
	var kingSquare Square
	for _, square := range board.Squares() {
		if !square.Empty && square.Piece.Identifier() == KingIdentifier && square.Piece.Color() == color {
			kingSquare = square
			break
		}
	}
	return kingSquare
}

func (g *game) oponentTurn() Player { // TODO: still needs test
	if g.Turn() == g.white {
		return g.black
	}
	return g.white
}
