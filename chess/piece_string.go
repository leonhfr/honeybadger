// Code generated by "stringer -type=Piece,PieceType,Color"; DO NOT EDIT.

package chess

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[WhitePawn-0]
	_ = x[WhiteKnight-2]
	_ = x[WhiteBishop-4]
	_ = x[WhiteRook-6]
	_ = x[WhiteQueen-8]
	_ = x[WhiteKing-10]
	_ = x[BlackPawn-1]
	_ = x[BlackKnight-3]
	_ = x[BlackBishop-5]
	_ = x[BlackRook-7]
	_ = x[BlackQueen-9]
	_ = x[BlackKing-11]
}

const _Piece_name = "WhitePawnBlackPawnWhiteKnightBlackKnightWhiteBishopBlackBishopWhiteRookBlackRookWhiteQueenBlackQueenWhiteKingBlackKing"

var _Piece_index = [...]uint8{0, 9, 18, 29, 40, 51, 62, 71, 80, 90, 100, 109, 118}

func (i Piece) String() string {
	if i >= Piece(len(_Piece_index)-1) {
		return "Piece(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Piece_name[_Piece_index[i]:_Piece_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Pawn-0]
	_ = x[Knight-2]
	_ = x[Bishop-4]
	_ = x[Rook-6]
	_ = x[Queen-8]
	_ = x[King-10]
}

const (
	_PieceType_name_0 = "Pawn"
	_PieceType_name_1 = "Knight"
	_PieceType_name_2 = "Bishop"
	_PieceType_name_3 = "Rook"
	_PieceType_name_4 = "Queen"
	_PieceType_name_5 = "King"
)

func (i PieceType) String() string {
	switch {
	case i == 0:
		return _PieceType_name_0
	case i == 2:
		return _PieceType_name_1
	case i == 4:
		return _PieceType_name_2
	case i == 6:
		return _PieceType_name_3
	case i == 8:
		return _PieceType_name_4
	case i == 10:
		return _PieceType_name_5
	default:
		return "PieceType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[White-0]
	_ = x[Black-1]
}

const _Color_name = "WhiteBlack"

var _Color_index = [...]uint8{0, 5, 10}

func (i Color) String() string {
	if i >= Color(len(_Color_index)-1) {
		return "Color(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Color_name[_Color_index[i]:_Color_index[i+1]]
}
