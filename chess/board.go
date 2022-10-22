package chess

import "strings"

// SquareMap represents a mapping between squares and pieces.
type SquareMap map[Square]Piece

// board represents a chess board and its relationship between squares and pieces.
type board struct {
	bbWhiteKing   bitboard
	bbWhiteQueen  bitboard
	bbWhiteRook   bitboard
	bbWhiteBishop bitboard
	bbWhiteKnight bitboard
	bbWhitePawn   bitboard
	bbBlackKing   bitboard
	bbBlackQueen  bitboard
	bbBlackRook   bitboard
	bbBlackBishop bitboard
	bbBlackKnight bitboard
	bbBlackPawn   bitboard
	bbWhite       bitboard
	bbBlack       bitboard
	bbEmpty       bitboard
	bbOccupied    bitboard
	sqWhiteKing   Square
	sqBlackKing   Square
}

func newBoard(m SquareMap) *board {
	b := &board{}
	for sq, p := range m {
		b.setPiece(p, sq)
	}
	b.computeConvenienceBitboards()
	return b
}

func (b *board) computeConvenienceBitboards() {
	b.bbWhite = b.bbWhiteKing | b.bbWhiteQueen | b.bbWhiteRook |
		b.bbWhiteBishop | b.bbWhiteKnight | b.bbWhitePawn
	b.bbBlack = b.bbBlackKing | b.bbBlackQueen | b.bbBlackRook |
		b.bbBlackBishop | b.bbBlackKnight | b.bbBlackPawn
	b.bbOccupied = b.bbWhite | b.bbBlack
	b.bbEmpty = ^b.bbOccupied

	for _, sq := range b.getBitboard(WhiteKing).mapping() {
		b.sqWhiteKing = sq
	}
	for _, sq := range b.getBitboard(BlackKing).mapping() {
		b.sqBlackKing = sq
	}
}

func (b *board) squareMap() SquareMap {
	m := SquareMap{}
	for _, p := range pieces {
		for _, sq := range b.getBitboard(p).mapping() {
			m[sq] = p
		}
	}
	return m
}

func (b *board) piece(sq Square) Piece {
	for _, p := range pieces {
		bb := b.getBitboard(p)
		if bb.occupied(sq) {
			return p
		}
	}
	return NoPiece
}

func (b *board) update(m Move) {
	p1, p2 := b.piece(m.S1()), b.piece(m.S2())
	// remove s1 piece
	b.removePiece(p1, m.S1())

	// remove s2 piece if any
	if p2 != NoPiece {
		b.removePiece(p2, m.S2())
	}

	// add s1 piece in s2
	b.setPiece(p1, m.S2())

	// promotion if any
	if m.Promo() != NoPiece {
		b.removePiece(newPiece(p1.Color(), Pawn), m.S2())
		b.setPiece(m.Promo(), m.S2())
	}

	// en passant
	switch c := p1.Color(); {
	case m.HasTag(EnPassant) && c == White:
		b.removePiece(BlackPawn, m.S2()-8)
	case m.HasTag(EnPassant) && c == Black:
		b.removePiece(WhitePawn, m.S2()+8)
	}

	// castle (only move rook)
	switch c := p1.Color(); {
	case c == White && m.HasTag(KingSideCastle):
		b.bbWhiteRook = b.bbWhiteRook & ^H1.bitboard() | F1.bitboard()
	case c == White && m.HasTag(QueenSideCastle):
		b.bbWhiteRook = b.bbWhiteRook & ^A1.bitboard() | D1.bitboard()
	case c == Black && m.HasTag(KingSideCastle):
		b.bbBlackRook = b.bbBlackRook & ^H8.bitboard() | F8.bitboard()
	case c == Black && m.HasTag(QueenSideCastle):
		b.bbBlackRook = b.bbBlackRook & ^A8.bitboard() | D8.bitboard()
	}

	b.computeConvenienceBitboards()
}

func (b *board) String() string {
	var fields []string
	for i := 7; i >= 0; i-- {
		r := Rank(8 * i)
		var field []byte
		for f := FileA; f <= FileH; f++ {
			if p := b.piece(NewSquare(f, r)); p != NoPiece {
				field = append(field, []byte(p.String())...)
			} else if len(field) == 0 {
				field = append(field, '1')
			} else if r := field[len(field)-1]; r < '1' || '8' < r {
				field = append(field, '1')
			} else {
				field[len(field)-1]++
			}
		}
		fields = append(fields, string(field))
	}
	return strings.Join(fields, "/")
}

func (b *board) setPiece(p Piece, sq Square) {
	switch bb := sq.bitboard(); p {
	case WhiteKing:
		b.bbWhiteKing |= bb
	case WhiteQueen:
		b.bbWhiteQueen |= bb
	case WhiteRook:
		b.bbWhiteRook |= bb
	case WhiteBishop:
		b.bbWhiteBishop |= bb
	case WhiteKnight:
		b.bbWhiteKnight |= bb
	case WhitePawn:
		b.bbWhitePawn |= bb
	case BlackKing:
		b.bbBlackKing |= bb
	case BlackQueen:
		b.bbBlackQueen |= bb
	case BlackRook:
		b.bbBlackRook |= bb
	case BlackBishop:
		b.bbBlackBishop |= bb
	case BlackKnight:
		b.bbBlackKnight |= bb
	case BlackPawn:
		b.bbBlackPawn |= bb
	}
}

func (b *board) removePiece(p Piece, sq Square) {
	bb := b.getBitboard(p) & ^sq.bitboard()
	b.setBitboard(p, bb)
}

func (b *board) getBitboard(p Piece) bitboard {
	switch p {
	case WhiteKing:
		return b.bbWhiteKing
	case WhiteQueen:
		return b.bbWhiteQueen
	case WhiteRook:
		return b.bbWhiteRook
	case WhiteBishop:
		return b.bbWhiteBishop
	case WhiteKnight:
		return b.bbWhiteKnight
	case WhitePawn:
		return b.bbWhitePawn
	case BlackKing:
		return b.bbBlackKing
	case BlackQueen:
		return b.bbBlackQueen
	case BlackRook:
		return b.bbBlackRook
	case BlackBishop:
		return b.bbBlackBishop
	case BlackKnight:
		return b.bbBlackKnight
	case BlackPawn:
		return b.bbBlackPawn
	default:
		panic("unknown piece")
	}
}

func (b *board) setBitboard(p Piece, bb bitboard) {
	switch p {
	case WhiteKing:
		b.bbWhiteKing = bb
	case WhiteQueen:
		b.bbWhiteQueen = bb
	case WhiteRook:
		b.bbWhiteRook = bb
	case WhiteBishop:
		b.bbWhiteBishop = bb
	case WhiteKnight:
		b.bbWhiteKnight = bb
	case WhitePawn:
		b.bbWhitePawn = bb
	case BlackKing:
		b.bbBlackKing = bb
	case BlackQueen:
		b.bbBlackQueen = bb
	case BlackRook:
		b.bbBlackRook = bb
	case BlackBishop:
		b.bbBlackBishop = bb
	case BlackKnight:
		b.bbBlackKnight = bb
	case BlackPawn:
		b.bbBlackPawn = bb
	}
}

func (b *board) copy() *board {
	return &board{
		bbWhiteKing:   b.bbWhiteKing,
		bbWhiteQueen:  b.bbWhiteQueen,
		bbWhiteRook:   b.bbWhiteRook,
		bbWhiteBishop: b.bbWhiteBishop,
		bbWhiteKnight: b.bbWhiteKnight,
		bbWhitePawn:   b.bbWhitePawn,
		bbBlackKing:   b.bbBlackKing,
		bbBlackQueen:  b.bbBlackQueen,
		bbBlackRook:   b.bbBlackRook,
		bbBlackBishop: b.bbBlackBishop,
		bbBlackKnight: b.bbBlackKnight,
		bbBlackPawn:   b.bbBlackPawn,
		bbWhite:       b.bbWhite,
		bbBlack:       b.bbBlack,
		bbEmpty:       b.bbEmpty,
	}
}
