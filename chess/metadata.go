package chess

// Metadata represents a position's metadata.
//
//	32 bits
//	__ square fullMove halfMove ___ CCCC T
//	square      en passant square
//	fullMove    full moves
//	halfMove    half move clock
//	CCCC        castle rights
//	T           turn color
//	_           unused bit
type Metadata uint32

func newMetadata(c Color, cr CastlingRights, halfMoveClock, fullMoves uint8, enPassant Square) Metadata {
	return Metadata(c) |
		Metadata(cr)<<1 |
		Metadata(halfMoveClock)<<8 |
		Metadata(fullMoves)<<16 |
		Metadata(enPassant)<<24
}

func (m Metadata) turn() Color {
	return Color(m & 1)
}

func (m Metadata) castleRights() CastlingRights {
	return CastlingRights((m >> 1) & 15)
}

func (m Metadata) halfMoveClock() uint8 {
	return uint8(m >> 8)
}

func (m Metadata) fullMoves() uint8 {
	return uint8(m >> 16)
}

func (m Metadata) enPassant() Square {
	return Square(m >> 24)
}

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
