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
	bbOccupied    bitboard
	bbPinned      bitboard // pinned pieces (can still move in direction of and attack pinner)
	bbPinner      bitboard // pieces that pin some opponent's pieces
	bbCheck       bitboard // pieces that give check
}

func newBoard(m SquareMap) board {
	b := board{}
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

	whiteKing, blackKing := b.getKingSquare(White), b.getKingSquare(Black)

	bbWhitePinned, bbWhitePinner := pinnedBitboard(whiteKing, b.bbOccupied,
		b.bbWhite, b.bbBlackQueen, b.bbBlackRook, b.bbBlackBishop)
	bbBlackPinned, bbBlackPinner := pinnedBitboard(blackKing, b.bbOccupied,
		b.bbBlack, b.bbWhiteQueen, b.bbWhiteRook, b.bbWhiteBishop)
	b.bbPinned = bbWhitePinned | bbBlackPinned
	b.bbPinner = bbWhitePinner | bbBlackPinner

	bbWhiteCheck := checkBitboard(whiteKing, White, b.bbOccupied,
		b.bbBlackKing, b.bbBlackQueen, b.bbBlackRook,
		b.bbBlackBishop, b.bbBlackKnight, b.bbBlackPawn)
	bbBlackCheck := checkBitboard(blackKing, Black, b.bbOccupied,
		b.bbWhiteKing, b.bbWhiteQueen, b.bbWhiteRook,
		b.bbWhiteBishop, b.bbWhiteKnight, b.bbWhitePawn)
	b.bbCheck = bbWhiteCheck | bbBlackCheck
}

func (b board) squareMap() SquareMap {
	m := SquareMap{}
	for p := BlackPawn; p <= WhiteKing; p++ {
		for _, sq := range b.getBitboard(p).mapping() {
			m[sq] = p
		}
	}
	return m
}

func (b board) piece(sq Square) Piece {
	for p := BlackPawn; p <= WhiteKing; p++ {
		if b.getBitboard(p).occupied(sq) {
			return p
		}
	}
	return NoPiece
}

func (b board) pieceByColor(sq Square, c Color) Piece {
	for p := newPiece(c, Pawn); p <= WhiteKing; p += 2 {
		if b.getBitboard(p).occupied(sq) {
			return p
		}
	}
	return NoPiece
}

func (b *board) makeMoveBoard(m Move) {
	p1, p2 := m.P1(), m.P2()
	s1, s2 := m.S1(), m.S2()

	// remove s1 piece
	b.removePiece(p1, s1)

	// remove s2 piece if any
	if p2 != NoPiece {
		b.removePiece(p2, s2)
	}

	// set s1 piece in s2
	// or set promotion piece if any
	if m.Promo() == NoPiece {
		b.setPiece(p1, s2)
	} else {
		b.setPiece(m.Promo(), s2)
	}

	// en passant and castle (only move rook)
	switch c := p1.Color(); {
	case m.HasTag(EnPassant) && c == White:
		b.removePiece(BlackPawn, s2-8)
	case m.HasTag(EnPassant) && c == Black:
		b.removePiece(WhitePawn, s2+8)
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

func (b *board) unmakeMoveBoard(m Move) {
	p1, p2 := m.P1(), m.P2()
	s1, s2 := m.S1(), m.S2()

	// remove original piece from new square
	// or remove promotion piece if any
	if m.Promo() == NoPiece {
		b.removePiece(p1, s2)
	} else {
		b.removePiece(m.Promo(), s2)
	}

	// set original piece in original square
	b.setPiece(p1, s1)

	// set captured piece if any
	if p2 != NoPiece {
		b.setPiece(p2, s2)
	}

	// en passant and castle (only move rook)
	switch c := p1.Color(); {
	case m.HasTag(EnPassant) && c == White:
		b.setPiece(BlackPawn, s2-8)
	case m.HasTag(EnPassant) && c == Black:
		b.setPiece(WhitePawn, s2+8)
	case c == White && m.HasTag(KingSideCastle):
		b.bbWhiteRook = b.bbWhiteRook & ^F1.bitboard() | H1.bitboard()
	case c == White && m.HasTag(QueenSideCastle):
		b.bbWhiteRook = b.bbWhiteRook & ^D1.bitboard() | A1.bitboard()
	case c == Black && m.HasTag(KingSideCastle):
		b.bbBlackRook = b.bbBlackRook & ^F8.bitboard() | H8.bitboard()
	case c == Black && m.HasTag(QueenSideCastle):
		b.bbBlackRook = b.bbBlackRook & ^D8.bitboard() | A8.bitboard()
	}

	b.computeConvenienceBitboards()
}

func (b board) hasSufficientMaterial() bool {
	if (b.bbWhiteQueen | b.bbWhiteRook | b.bbWhitePawn |
		b.bbBlackQueen | b.bbBlackRook | b.bbBlackPawn) > 0 {
		return true
	}

	nWhites, nBlacks := b.bbWhite.ones(), b.bbBlack.ones()
	nWhiteBishops, nBlackBishops := b.bbWhiteBishop.ones(), b.bbBlackBishop.ones()

	// king versus king
	// king and bishop versus king
	// king and knight versus king
	if b.bbWhite == b.bbWhiteKing && b.bbBlack == b.bbBlackKing ||
		b.bbWhite == b.bbWhiteKing && nBlacks == 2 && (nBlackBishops == 1 || b.bbBlackKnight.ones() == 1) ||
		b.bbBlack == b.bbBlackKing && nWhites == 2 && (nWhiteBishops == 1 || b.bbWhiteKnight.ones() == 1) {
		return false
	}

	// king and bishop versus king and bishop with the bishops on the same color
	if nWhites == 2 && nBlacks == 2 && nWhiteBishops == 1 && nBlackBishops == 1 &&
		(((b.bbBlackBishop|b.bbWhiteBishop)&bbWhiteSquares).ones() == 2 ||
			((b.bbBlackBishop|b.bbWhiteBishop)&bbBlackSquares).ones() == 2) {
		return false
	}

	return true
}

func (b board) String() string {
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
	switch mask := ^sq.bitboard(); p {
	case WhiteKing:
		b.bbWhiteKing &= mask
	case WhiteQueen:
		b.bbWhiteQueen &= mask
	case WhiteRook:
		b.bbWhiteRook &= mask
	case WhiteBishop:
		b.bbWhiteBishop &= mask
	case WhiteKnight:
		b.bbWhiteKnight &= mask
	case WhitePawn:
		b.bbWhitePawn &= mask
	case BlackKing:
		b.bbBlackKing &= mask
	case BlackQueen:
		b.bbBlackQueen &= mask
	case BlackRook:
		b.bbBlackRook &= mask
	case BlackBishop:
		b.bbBlackBishop &= mask
	case BlackKnight:
		b.bbBlackKnight &= mask
	case BlackPawn:
		b.bbBlackPawn &= mask
	}
}

func (b board) getBitboard(p Piece) bitboard {
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

func (b board) getColor(c Color) bitboard {
	if c == White {
		return b.bbWhite
	}
	return b.bbBlack
}

func (b board) getPinned(c Color) bitboard {
	if c == White {
		return b.bbPinned & b.bbWhite
	}
	return b.bbPinned & b.bbBlack
}

func (b board) getPinner(c Color) bitboard {
	if c == White {
		return b.bbPinner & b.bbWhite
	}
	return b.bbPinner & b.bbBlack
}

func (b board) getCheck(c Color) bitboard {
	if c == White {
		return b.bbCheck & b.bbWhite
	}
	return b.bbCheck & b.bbBlack
}

func (b board) getKingSquare(c Color) Square {
	if c == White {
		return b.bbWhiteKing.scanForward()
	}
	return b.bbBlackKing.scanForward()
}

func (b board) copyBoard() board {
	return board{
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
		bbOccupied:    b.bbOccupied,
		bbPinned:      b.bbPinned,
		bbPinner:      b.bbPinner,
		bbCheck:       b.bbCheck,
	}
}
