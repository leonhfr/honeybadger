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
	b.bbEmpty = ^(b.bbWhite | b.bbBlack)
}

func (b *board) squareMap() SquareMap {
	m := SquareMap{}
	for _, p := range pieces {
		squares := b.getBitboard(p).mapping()
		for sq := range squares {
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

func (b *board) update(m *Move) {
	p1, p2 := b.piece(m.s1), b.piece(m.s2)
	// remove s1 piece
	b.removePiece(p1, m.s1)

	// remove s2 piece if any
	if p2 != NoPiece {
		b.removePiece(p2, m.s2)
	}

	// add s1 piece in s2
	b.setPiece(p1, m.s2)

	// promotion if any
	if m.promo != NoPieceType {
		b.removePiece(newPiece(p1.Color(), Pawn), m.s2)
		b.setPiece(newPiece(p1.Color(), m.promo), m.s2)
	}

	// en passant
	switch c := p1.Color(); {
	case m.HasTag(EnPassant) && c == White:
		b.removePiece(BlackPawn, m.s2-8)
	case m.HasTag(EnPassant) && c == Black:
		b.removePiece(WhitePawn, m.s2+8)
	}

	// castle (only move rook)
	switch c := p1.Color(); {
	case m.HasTag(KingSideCastle) && c == White:
		b.bbWhiteRook = b.bbWhiteRook & ^bbForSquare(H1) | bbForSquare(F1)
	case m.HasTag(QueenSideCastle) && c == White:
		b.bbWhiteRook = b.bbWhiteRook & ^bbForSquare(A1) | bbForSquare(D1)
	case m.HasTag(KingSideCastle) && c == Black:
		b.bbBlackRook = b.bbBlackRook & ^bbForSquare(H8) | bbForSquare(F8)
	case m.HasTag(QueenSideCastle) && c == Black:
		b.bbBlackRook = b.bbBlackRook & ^bbForSquare(A8) | bbForSquare(D8)
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
	switch bb := bbForSquare(sq); p {
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
	bb := b.getBitboard(p) & ^bbForSquare(sq)
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

func bbForSquare(sq Square) bitboard {
	return 1 << sq
}
