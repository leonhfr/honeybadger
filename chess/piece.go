package chess

var (
	pieces = [12]Piece{
		WhiteKing, WhiteQueen, WhiteRook, WhiteBishop, WhiteKnight, WhitePawn,
		BlackKing, BlackQueen, BlackRook, BlackBishop, BlackKnight, BlackPawn,
	}
	piecesByColor = [2][6]Piece{
		{WhiteKing, WhiteQueen, WhiteRook, WhiteBishop, WhiteKnight, WhitePawn},
		{BlackKing, BlackQueen, BlackRook, BlackBishop, BlackKnight, BlackPawn},
	}
)

// Color represents the color of a chess piece.
type Color uint8

const (
	// White represents the white color.
	White Color = iota
	// Black represents the black color.
	Black
)

// Other returns the other color.
func (c Color) Other() Color {
	if c == White {
		return Black
	}
	return White
}

const colorName = "wb"

// String implements the Stringer interface.
// Returns an UCI-compatible representation.
func (c Color) String() string {
	return colorName[c : c+1]
}

// PieceType is the type of a piece.
type PieceType uint8

const (
	// Pawn represents a pawn.
	Pawn PieceType = iota << 1
	// Knight represents a knight.
	Knight
	// Bishop represents a bishop.
	Bishop
	// Rook represents a rook.
	Rook
	// Queen represents a queen.
	Queen
	// King represents a king.
	King
	// NoPieceType represents an absence of PieceType.
	NoPieceType
)

const pieceTypeName = "pnbrqk-"

// String implements the Stringer interface.
// Returns an UCI-compatible representation.
func (pt PieceType) String() string {
	return pieceTypeName[pt/2 : pt/2+1]
}

// Piece is a piece type with a color.
type Piece uint8

const (
	// WhitePawn represents a white pawn.
	WhitePawn Piece = Piece(White) | Piece(Pawn)
	// WhiteKnight represents a white knight.
	WhiteKnight Piece = Piece(White) | Piece(Knight)
	// WhiteBishop represents a white bishop.
	WhiteBishop Piece = Piece(White) | Piece(Bishop)
	// WhiteRook represents a white rook.
	WhiteRook Piece = Piece(White) | Piece(Rook)
	// WhiteQueen represents a white queen.
	WhiteQueen Piece = Piece(White) | Piece(Queen)
	// WhiteKing represents a white king.
	WhiteKing Piece = Piece(White) | Piece(King)
	// BlackPawn represents a black pawn.
	BlackPawn Piece = Piece(Black) | Piece(Pawn)
	// BlackKnight represents a black knight.
	BlackKnight Piece = Piece(Black) | Piece(Knight)
	// BlackBishop represents a black bishop.
	BlackBishop Piece = Piece(Black) | Piece(Bishop)
	// BlackRook represents a black rook.
	BlackRook Piece = Piece(Black) | Piece(Rook)
	// BlackQueen represents a black queen.
	BlackQueen Piece = Piece(Black) | Piece(Queen)
	// BlackKing represents a black king.
	BlackKing Piece = Piece(Black) | Piece(King)
	// NoPiece represents an absence of Piece.
	NoPiece Piece = 12
)

func newPiece(c Color, pt PieceType) Piece {
	return Piece(c) | Piece(pt)
}

const pieceName = "PpNnBbRrQqKk-"

// String implements the Stringer interface.
// Returns an UCI-compatible representation.
func (p Piece) String() string {
	return pieceName[p : p+1]
}

// Color returns the color of the piece.
func (p Piece) Color() Color {
	return Color(p & 1)
}

// Type returns the type of the piece.
func (p Piece) Type() PieceType {
	return PieceType(p & ^Piece(1))
}
