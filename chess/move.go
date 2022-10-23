package chess

import "errors"

var (
	errInvalidMove     = errors.New("invalid move in UCI notation")
	errMissingPosition = errors.New("missing position")
)

// MoveTag represents a notable consequence of a move.
type MoveTag uint32

const (
	// Quiet indicates that the move is a priori quiet.
	Quiet MoveTag = 1 << (iota + 24)
	// Capture indicates that the move captures a piece.
	Capture
	// EnPassant indicates that the move captures a piece via en passant.
	EnPassant
	// Promotion indicates that the move is a promotion.
	Promotion
	// KingSideCastle indicates that the move is a king side castle.
	KingSideCastle
	// QueenSideCastle indicates that the move is a queen side castle.
	QueenSideCastle
)

// Move represents a move from a square to another.
//
//	32 bits
//	xxxxxxxx pppp tttt ffff TTTTTT FFFFFF
//	xxxxxxxx   move tags
//	pppp       promo piece
//	tttt       to piece
//	ffff       from piece
//	TTTTTT     to square
//	FFFFFF     from square
type Move uint32

func newMove(p1, p2 Piece, s1, s2, enPassant Square, promo Piece) Move {
	var tags MoveTag
	if pt := p1.Type(); pt == King {
		if (s1 == E1 && s2 == G1) || (s1 == E8 && s2 == G8) {
			tags |= KingSideCastle
		} else if (s1 == E1 && s2 == C1) || (s1 == E8 && s2 == C8) {
			tags |= QueenSideCastle
		}
	} else if pt == Pawn && s2 == enPassant {
		tags |= EnPassant
		tags |= Capture
	} else if promo != NoPiece {
		tags |= Promotion
	}

	if p2 != NoPiece {
		tags |= Capture
	}

	if tags == 0 {
		tags |= Quiet
	}

	return Move(s1) | Move(s2)<<6 |
		Move(p1)<<12 | Move(p2&15)<<16 |
		Move(promo&15)<<20 | Move(tags)
}

// S1 returns the origin square of the move.
func (m Move) S1() Square {
	return Square(m & 63)
}

// S2 returns the destination square of the move.
func (m Move) S2() Square {
	return Square((m >> 6) & 63)
}

// P1 returns the piece in the origin square.
func (m Move) P1() Piece {
	return Piece((m >> 12) & 15)
}

// P2 returns the piece in the destination square.
func (m Move) P2() Piece {
	return Piece((m >> 16) & 15)
}

// Promo returns the promotion piece of the move.
func (m Move) Promo() Piece {
	return Piece((m >> 20) & 15)
}

// HasTag checks whether the move has the given MoveTag.
func (m Move) HasTag(tag MoveTag) bool {
	return tag&MoveTag(m) > 0
}

// String implements the Stringer interface.
// Returns the move in UCI notation.
func (m Move) String() string {
	base := m.S1().String() + m.S2().String()
	if promo := m.Promo(); promo != NoPiece {
		base += promo.Type().String()
	}
	return base
}

// MoveFromUCI creates a move from a string in UCI notation.
func MoveFromUCI(pos *Position, s string) (Move, error) {
	if pos == nil {
		return 0, errMissingPosition
	}

	if len(s) < 4 || len(s) > 5 {
		return 0, errInvalidMove
	}

	s1, ok := strToSquareMap[s[0:2]]
	if !ok {
		return 0, errInvalidMove
	}

	s2, ok := strToSquareMap[s[2:4]]
	if !ok {
		return 0, errInvalidMove
	}

	promo := NoPiece
	if len(s) == 5 {
		promoType, ok := uciPieceTypeMap[s[4:5]]
		promo = newPiece(pos.turn, promoType)
		if !ok {
			return 0, errInvalidMove
		}
	}

	p1 := pos.board.piece(s1)
	p2 := pos.board.piece(s2)
	return newMove(p1, p2, s1, s2, pos.enPassant, promo), nil
}

var uciPieceTypeMap = map[string]PieceType{
	"q": Queen,
	"r": Rook,
	"b": Bishop,
	"n": Knight,
}
