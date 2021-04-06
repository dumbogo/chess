package engine

// Game current game
type Game struct {
	Name  string
	Board Board
	turn  Player
}

// Player gammer
type Player struct {
	Name string
}

// Move moves a piece in the Board
func (b *Board) Move(from, to Square) bool {
	// TODO: WIP
	var isPosible bool
	// Get piece
	// Review if piece can make the movement
	// Review what "to" square has
	//	1. nothing, pass
	//	2. has piece:
	//		1. from == to, not pass
	//		1. from != to, eat and pass

	return isPosible
}

// RemovePiece removes a piece from Board location
func (b *Board) RemovePiece(loc SquareIdentifier) Piece {
	// TODO: WIP
	// SHould remove piece game logic be located here ?

	var piece Piece
	square := b.Squares[loc]
	square.Empty = true
	piece = square.Piece
	// 	square.Piece = Piece{}
	return piece
}
