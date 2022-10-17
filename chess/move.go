package chess

import "errors"

var errInvalidMove = errors.New("invalid move in UCI notation")

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
)

// Move represents a move from a square to another.
type Move struct {
	s1    Square
	s2    Square
	promo PieceType
	tags  MoveTag
}

// FromUCI creates a move from a string in UCI notation.
func FromUCI(p *Position, s string) (*Move, error) {
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

	m := &Move{s1, s2, promo, 0}
	if p == nil {
		return m, nil
	}

	m.tags = moveTags(p, m)
	return m, nil
}

func moveTags(p *Position, m *Move) MoveTag {
	var tags MoveTag
	p1, p2 := p.board.piece(m.s1), p.board.piece(m.s2)
	c1, c2 := p1.Color(), p2.Color()
	t1 := p1.Type()

	if t1 == King {
		if (m.s1 == E1 && m.s2 == G1) || (m.s1 == E8 && m.s2 == G8) {
			tags |= KingSideCastle
		} else if (m.s1 == E1 && m.s2 == C1) || (m.s1 == E8 && m.s2 == C8) {
			tags |= QueenSideCastle
		}
	} else if t1 == Pawn && m.s2 == p.enPassantSquare {
		tags |= EnPassant
		tags |= Capture
	}

	if p2 != NoPiece && c1 != c2 {
		tags |= Capture
	}

	return tags
}

// String implements the Stringer interface.
// Returns the move in UCI notation.
func (m Move) String() string {
	base := m.s1.String() + m.S2().String()
	if m.promo != Pawn {
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
