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

// Method is the method that generated the outcome.
type Method uint8

const (
	// NoMethod indicates that an outcome hasn't occurred or that the method can't be determined.
	NoMethod Method = iota
	// Checkmate indicates that the game was won checkmate.
	Checkmate
	// Stalemate indicates that the game was drawn by stalemate.
	// The player whose turn it is to move is not in check and has no legal move.
	Stalemate
	// InsufficientMaterial indicates that the game was automatically drawn
	// because there was insufficient material for checkmate.
	InsufficientMaterial
)

// Outcome is the result of a game.
type Outcome uint8

const (
	// NoOutcome indicates that a game is in progress or ended without a result.
	NoOutcome Outcome = iota
	// WhiteWon indicates that white won the game.
	WhiteWon
	// BlackWon indicates that black won the game.
	BlackWon
	// Draw indicates that game was a draw.
	Draw
)
