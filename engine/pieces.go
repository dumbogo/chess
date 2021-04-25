package engine

// Piece piece type
type Piece interface {
	// Name Piece name
	Identifier() PieceIdentifier
	// Color piece color
	Color() Color
	// CanMove check if piece can be moved from, to Square,
	// returns true if it's possible, even if a piece can be eaten
	CanMove(board Board, movements []Movement, from, to Square) bool

	String() string
}

type piece struct {
	pieceIdentifier PieceIdentifier
	color           Color
}

// Name interface Piece
func (p *piece) Identifier() PieceIdentifier {
	return p.pieceIdentifier
}

func (p *piece) Color() Color {
	return p.color
}

type pawn struct {
	piece
}

// NewPawn creates Pawn piece
func NewPawn(color Color) Piece {
	p := piece{
		pieceIdentifier: PawnIdentifier,
		color:           color,
	}
	return &pawn{
		piece: p,
	}
}

// TODO: PawnCanMove [En passant , Handle when reachs Coordinates limits]
func (p *pawn) CanMove(board Board, movements []Movement, from, to Square) bool {
	if from.Empty {
		return false
	}
	if from.Piece.Identifier() != PawnIdentifier {
		return false
	}

	if from.Piece.Color() == WhiteColor {
		// First movement , not eating
		if from.Coordinates.Y == 1 &&
			from.Coordinates.X == to.Coordinates.X &&
			to.Coordinates.Y >= from.Coordinates.Y+2 &&
			to.Empty {
			// BUG : when there's a piece in position X+1 and pawn moves two spaces, it returns true
			return true
		}

		// Simple movement, advance one
		if to.Empty &&
			from.Coordinates.Y+1 == to.Coordinates.Y &&
			from.Coordinates.X == to.Coordinates.X {
			return true
		}

		// Eating piece
		if from.Coordinates.Y+1 == to.Coordinates.Y &&
			(from.Coordinates.X-1 == to.Coordinates.X || from.Coordinates.X+1 == to.Coordinates.X) &&
			!to.Empty &&
			to.Piece.Color() != from.Piece.Color() {
			return true
		}
	} else if from.Piece.Color() == BlackColor {
		// First movement , not eating
		if from.Coordinates.Y == 6 &&
			from.Coordinates.X == to.Coordinates.X &&
			to.Coordinates.Y >= from.Coordinates.Y-2 &&
			to.Empty {
			// BUG : when there's a piece in position X-1 and pawn moves two spaces, it returns true
			return true
		}

		// Simple movement, advance one
		if to.Empty &&
			from.Coordinates.Y-1 == to.Coordinates.Y &&
			from.Coordinates.X == to.Coordinates.X {
			return true
		}

		// Eating piece
		if from.Coordinates.Y-1 == to.Coordinates.Y &&
			(from.Coordinates.X-1 == to.Coordinates.X || from.Coordinates.X+1 == to.Coordinates.X) &&
			!to.Empty &&
			to.Piece.Color() != from.Piece.Color() {
			return true
		}
	}
	return false
}

func (p *pawn) String() string {
	if p.Color() == WhiteColor {
		return "WP"
	}
	return "BP"
}

type king struct {
	piece
}

// NewKing creates King piece
func NewKing(color Color) Piece {
	p := piece{
		pieceIdentifier: KingIdentifier,
		color:           color,
	}
	return &king{
		piece: p,
	}
}

func (k *king) CanMove(b Board, m []Movement, from, to Square) bool {
	if from.Empty {
		return false
	}
	if from.Piece.Identifier() != KingIdentifier {
		return false
	}

	for _, square := range b.Squares() {
		if !square.Empty && square.Piece.Color() != k.Color() {
			if square.Piece.CanMove(b, m, square, to) {
				return false
			}
		}
	}

	if (to.Coordinates.X == from.Coordinates.X+1 || to.Coordinates.X == from.Coordinates.X-1) &&
		(to.Coordinates.Y == from.Coordinates.Y+1 || to.Coordinates.Y == from.Coordinates.Y-1) &&
		(to.Empty || !to.Empty && to.Piece.Color() != from.Piece.Color()) {
		return true
	}

	return false
}

func (k *king) String() string {
	if k.Color() == WhiteColor {
		return "WK"
	}
	return "BK"
}

type queen struct {
	piece
}

// NewQueen creates Queen piece
func NewQueen(color Color) Piece {
	p := piece{
		pieceIdentifier: QueenIdentifier,
		color:           color,
	}
	return &queen{
		piece: p,
	}
}

func (q *queen) CanMove(b Board, m []Movement, from, to Square) bool {
	// TODO: CanMove queen
	panic("TODO: ENDTHIS")
}

func (q *queen) String() string {
	if q.Color() == WhiteColor {
		return "WQ"
	}
	return "BQ"
}

type rook struct {
	piece
}

// NewRook creates Rook piece
func NewRook(color Color) Piece {
	p := piece{
		pieceIdentifier: RookIdentifier,
		color:           color,
	}
	return &rook{
		piece: p,
	}
}

func (r *rook) CanMove(b Board, m []Movement, from, to Square) bool {
	if from.Empty {
		return false
	}
	squares := b.Squares()
	if from.Coordinates.X == to.Coordinates.X &&
		from.Coordinates.Y != to.Coordinates.Y {

		if !to.Empty && to.Piece.Color() == from.Piece.Color() {
			return false
		}

		// we need to make sure that, for each square between from and to, are empty
		summ := 1
		if from.Coordinates.Y > to.Coordinates.Y {
			summ = -1
		}
		for Yiterator := (int(from.Coordinates.Y) + summ); Yiterator != int(to.Coordinates.Y); Yiterator += summ {
			square := squares[CoordinateToSquareIdentifier(Coordinate{X: from.Coordinates.X, Y: uint8(Yiterator)})]
			if !square.Empty {
				return false
			}
		}
		return true
	} else if from.Coordinates.Y == to.Coordinates.Y && from.Coordinates.X != to.Coordinates.X {
		if !to.Empty && to.Piece.Color() == from.Piece.Color() {
			return false
		}

		// we need to make sure that, for each square between from and to, are empty
		summ := 1
		if from.Coordinates.X > to.Coordinates.X {
			summ = -1
		}
		for Xiterator := (int(from.Coordinates.X) + summ); Xiterator != int(to.Coordinates.X); Xiterator += summ {
			square := squares[CoordinateToSquareIdentifier(Coordinate{X: uint8(Xiterator), Y: from.Coordinates.Y})]
			if !square.Empty {
				return false
			}
		}
		return true

	}
	return false
}

func (r *rook) String() string {
	if r.Color() == WhiteColor {
		return "WR"
	}
	return "BR"
}

type bishop struct {
	piece
}

// NewBishop creates Bishop piece
func NewBishop(color Color) Piece {
	p := piece{
		pieceIdentifier: BishopIdentifier,
		color:           color,
	}
	return &bishop{
		piece: p,
	}
}

func (bi *bishop) CanMove(b Board, m []Movement, from, to Square) bool {
	// TODO: CanMove bishop
	panic("TODO: ENDTHIS")
}

func (bi *bishop) String() string {
	if bi.Color() == WhiteColor {
		return "WB"
	}
	return "BB"
}

type knight struct {
	piece
}

// NewKnight creates knight piece
func NewKnight(color Color) Piece {
	p := piece{
		pieceIdentifier: KnightIdentifier,
		color:           color,
	}
	return &knight{
		piece: p,
	}
}

func (k *knight) CanMove(b Board, m []Movement, from, to Square) bool {
	// TODO: CanMove knight
	panic("TODO: ENDTHIS")
}

func (k *knight) String() string {
	if k.Color() == WhiteColor {
		return "Wk"
	}
	return "Bk"
}
