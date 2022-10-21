package chess

import "errors"

var (
	errIllegalMove     = errors.New("illegal move")
	errInvalidMove     = errors.New("invalid move in UCI notation")
	errMissingPosition = errors.New("missing position")
)

// MoveTag represents a notable consequence of a move.
type MoveTag uint8

const (
	// KingSideCastle indicates that the move is a king side castle.
	KingSideCastle MoveTag = 1 << iota
	// QueenSideCastle indicates that the move is a queen side castle.
	QueenSideCastle
	// Capture indicates that the move captures a piece.
	Capture
	// EnPassant indicates that the move captures a piece via en passant.
	EnPassant
	// Check indicates that the move puts the opposing player in check.
	Check
	// inCheck indicates the the move puts the moving player in check. Illegal move.
	inCheck
)

// Move represents a move from a square to another.
type Move struct {
	s1    Square
	s2    Square
	promo PieceType
	tags  MoveTag
}

func newMove(pos *Position, pt PieceType, s1, s2 Square, promo PieceType) *Move {
	tags := moveTags(pos, pt, s1, s2)
	m := &Move{s1, s2, promo, tags}

	next := pos.Move(m)
	if isInCheck(next) {
		m.tags |= Check
	}
	next.turn = next.turn.Other()
	if isInCheck(next) {
		m.tags |= inCheck
	}

	return m
}

// FromUCI creates a move from a string in UCI notation.
func FromUCI(pos *Position, s string) (*Move, error) {
	if pos == nil {
		return nil, errMissingPosition
	}

	if len(s) < 4 || len(s) > 5 {
		return nil, errInvalidMove
	}

	s1, ok := strToSquareMap[s[0:2]]
	if !ok {
		return nil, errInvalidMove
	}

	s2, ok := strToSquareMap[s[2:4]]
	if !ok {
		return nil, errInvalidMove
	}

	promo := NoPieceType
	if len(s) == 5 {
		promo, ok = uciPieceTypeMap[s[4:5]]
		if !ok {
			return nil, errInvalidMove
		}
	}

	m := newMove(pos, NoPieceType, s1, s2, promo)
	if m.HasTag(inCheck) {
		return nil, errIllegalMove
	}

	return m, nil
}

func moveTags(pos *Position, pt PieceType, s1, s2 Square) MoveTag {
	if pt == NoPieceType {
		pt = pos.board.piece(s1).Type()
	}

	var tags MoveTag
	p2 := pos.board.piece(s2)

	if pt == King {
		if (s1 == E1 && s2 == G1) || (s1 == E8 && s2 == G8) {
			tags |= KingSideCastle
		} else if (s1 == E1 && s2 == C1) || (s1 == E8 && s2 == C8) {
			tags |= QueenSideCastle
		}
	} else if pt == Pawn && s2 == pos.enPassantSquare {
		tags |= EnPassant
		tags |= Capture
	}

	if p2 != NoPiece {
		tags |= Capture
	}

	return tags
}

// String implements the Stringer interface.
// Returns the move in UCI notation.
func (m Move) String() string {
	base := m.s1.String() + m.S2().String()
	if m.promo != NoPieceType {
		base += m.promo.String()
	}
	return base
}

// S1 returns the origin square of the move.
func (m Move) S1() Square {
	return m.s1
}

// S2 returns the destination square of the move.
func (m Move) S2() Square {
	return m.s2
}

// Promo returns the promotion piece type of the move.
func (m Move) Promo() PieceType {
	return m.promo
}

// HasTag checks whether the move has the given MoveTag.
func (m Move) HasTag(tag MoveTag) bool {
	return (tag & m.tags) > 0
}

var uciPieceTypeMap = map[string]PieceType{
	"q": Queen,
	"r": Rook,
	"b": Bishop,
	"n": Knight,
}
