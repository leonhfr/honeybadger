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
		b.setPiece(sq, p)
	}
	b.computeConvenienceBitboards()
	return b
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

func (b *board) computeConvenienceBitboards() {
	b.bbWhite = b.bbWhiteKing | b.bbWhiteQueen | b.bbWhiteRook |
		b.bbWhiteBishop | b.bbWhiteKnight | b.bbWhitePawn
	b.bbBlack = b.bbBlackKing | b.bbBlackQueen | b.bbBlackRook |
		b.bbBlackBishop | b.bbBlackKnight | b.bbBlackPawn
	b.bbEmpty = ^(b.bbWhite | b.bbBlack)
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

func (b *board) setPiece(sq Square, p Piece) {
	switch p {
	case WhiteKing:
		b.bbWhiteKing |= 1 << sq
	case WhiteQueen:
		b.bbWhiteQueen |= 1 << sq
	case WhiteRook:
		b.bbWhiteRook |= 1 << sq
	case WhiteBishop:
		b.bbWhiteBishop |= 1 << sq
	case WhiteKnight:
		b.bbWhiteKnight |= 1 << sq
	case WhitePawn:
		b.bbWhitePawn |= 1 << sq
	case BlackKing:
		b.bbBlackKing |= 1 << sq
	case BlackQueen:
		b.bbBlackQueen |= 1 << sq
	case BlackRook:
		b.bbBlackRook |= 1 << sq
	case BlackBishop:
		b.bbBlackBishop |= 1 << sq
	case BlackKnight:
		b.bbBlackKnight |= 1 << sq
	case BlackPawn:
		b.bbBlackPawn |= 1 << sq
	}
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
