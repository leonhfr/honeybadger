package chess

// Side represents a side of the board.
type Side uint8

const (
	// KingSide represents the kings' side.
	KingSide Side = iota
	// QueenSide represents the queens' side.
	QueenSide
)

// CastlingRights represents the castling right of one combination of side and color.
type CastlingRights uint8

const (
	// CastleWhiteKing represents white's king castle.
	CastleWhiteKing CastlingRights = 1 << iota
	// CastleWhiteQueen represents white's queen castle.
	CastleWhiteQueen
	// CastleBlackKing represents black's king castle.
	CastleBlackKing
	// CastleBlackQueen represents black's queen castle.
	CastleBlackQueen
)

// CanCastle returns whether a castle with this combinations of
// color and side is possible.
func (cr CastlingRights) CanCastle(c Color, s Side) bool {
	switch {
	case c == White && s == KingSide:
		return (cr & CastleWhiteKing) > 0
	case c == White && s == QueenSide:
		return (cr & CastleWhiteQueen) > 0
	case c == Black && s == KingSide:
		return (cr & CastleBlackKing) > 0
	case c == Black && s == QueenSide:
		return (cr & CastleBlackQueen) > 0
	default:
		return false
	}
}

func (cr CastlingRights) String() string {
	var rights string
	if (cr & CastleWhiteKing) > 0 {
		rights += "K"
	}
	if (cr & CastleWhiteQueen) > 0 {
		rights += "Q"
	}
	if (cr & CastleBlackKing) > 0 {
		rights += "k"
	}
	if (cr & CastleBlackQueen) > 0 {
		rights += "q"
	}
	if len(rights) == 0 {
		rights += "-"
	}
	return rights
}
