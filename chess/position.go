package chess

// CastlingRight represents the castling right of one combination of side and color.
type CastlingRight uint8

const (
	// CastleWhiteKing represents white's king castle.
	CastleWhiteKing CastlingRight = 1 << iota
	// CastleWhiteQueen represents white's queen castle.
	CastleWhiteQueen
	// CastleBlackKing represents black's king castle.
	CastleBlackKing
	// CastleBlackQueen represents black's queen castle.
	CastleBlackQueen
)
