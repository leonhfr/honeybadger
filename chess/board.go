package chess

import "strings"

// SquareMap represents a mapping between squares and pieces.
type SquareMap map[Square]Piece

// board represents a chess board and its relationship between squares and pieces.
type board struct {
	bbKing     bitboard
	bbQueen    bitboard
	bbRook     bitboard
	bbBishop   bitboard
	bbKnight   bitboard
	bbPawn     bitboard
	bbWhite    bitboard
	bbBlack    bitboard
	bbOccupied bitboard
	bbPinned   bitboard // pinned pieces (can still move in direction of and attack pinner)
	bbPinner   bitboard // pieces that pin some opponent's pieces
	bbCheck    bitboard // pieces that give check
}

func newBoard(m SquareMap) board {
	b := board{}
	for sq, p := range m {
		bb := sq.bitboard()
		b.xorBitboard(p.Type(), bb)
		b.xorColor(p.Color(), bb)
		b.bbOccupied ^= bb
	}
	b.computeConvenienceBitboards()
	return b
}

func (b *board) computeConvenienceBitboards() {
	whiteKing, blackKing := b.getKingSquare(White), b.getKingSquare(Black)

	bbWhitePinned, bbWhitePinner := pinnedBitboard(whiteKing, b.bbOccupied,
		b.bbWhite, b.bbBlack&b.bbQueen, b.bbBlack&b.bbRook, b.bbBlack&b.bbBishop)
	bbBlackPinned, bbBlackPinner := pinnedBitboard(blackKing, b.bbOccupied,
		b.bbBlack, b.bbWhite&b.bbQueen, b.bbWhite&b.bbRook, b.bbWhite&b.bbBishop)
	b.bbPinned = bbWhitePinned | bbBlackPinned
	b.bbPinner = bbWhitePinner | bbBlackPinner

	bbWhiteCheck := checkBitboard(whiteKing, White, b.bbOccupied,
		b.bbBlack&b.bbKing, b.bbBlack&b.bbQueen, b.bbBlack&b.bbRook,
		b.bbBlack&b.bbBishop, b.bbBlack&b.bbKnight, b.bbBlack&b.bbPawn)
	bbBlackCheck := checkBitboard(blackKing, Black, b.bbOccupied,
		b.bbWhite&b.bbKing, b.bbWhite&b.bbQueen, b.bbWhite&b.bbRook,
		b.bbWhite&b.bbBishop, b.bbWhite&b.bbKnight, b.bbWhite&b.bbPawn)
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

func (b board) pieceAt(sq Square) Piece {
	p := BlackPawn
	if b.bbWhite&sq.bitboard() > 0 {
		p = WhitePawn
	}
	for ; p <= WhiteKing; p += 2 {
		if b.getBitboard(p).occupied(sq) {
			return p
		}
	}
	return NoPiece
}

func (b board) pieceByColor(sq Square, c Color) Piece {
	for p := Pawn.color(c); p <= WhiteKing; p += 2 {
		if b.getBitboard(p).occupied(sq) {
			return p
		}
	}
	return NoPiece
}

func (b *board) makeMoveBoard(m Move) {
	p1, p2 := m.P1(), m.P2()
	s1, s2 := m.S1(), m.S2()
	c := p1.Color()

	s1bb, s2bb := s1.bitboard(), s2.bitboard()
	mbb := s1bb ^ s2bb

	if promo := m.Promo(); promo == NoPiece {
		b.xorBitboard(p1.Type(), mbb)
	} else {
		// promotion
		b.xorBitboard(p1.Type(), s1bb)
		b.xorBitboard(promo.Type(), s2bb)
	}

	b.xorColor(c, mbb)

	switch enPassant := m.HasTag(EnPassant); {
	case m.HasTag(Capture) && !enPassant: // capture
		b.xorBitboard(p2.Type(), s2bb)
		b.xorColor(p2.Color(), s2bb)
		b.bbOccupied ^= s1bb
	case c == White && enPassant: // white en passant
		bb := (s2 - 8).bitboard()
		b.bbPawn ^= bb
		b.bbBlack ^= bb
		b.bbOccupied ^= mbb ^ bb
	case c == Black && enPassant: // black en passant
		bb := (s2 + 8).bitboard()
		b.bbPawn ^= bb
		b.bbWhite ^= bb
		b.bbOccupied ^= mbb ^ bb
	case c == White && m.HasTag(KingSideCastle): // white king side castle
		b.bbRook ^= bbWhiteKingCastle
		b.bbWhite ^= bbWhiteKingCastle
		b.bbOccupied ^= bbWhiteKingCastleTravel
	case c == White && m.HasTag(QueenSideCastle): // white queen side castle
		b.bbRook ^= bbWhiteQueenCastle
		b.bbWhite ^= bbWhiteQueenCastle
		b.bbOccupied ^= bbWhiteQueenCastleTravel
	case c == Black && m.HasTag(KingSideCastle): // black king side castle
		b.bbRook ^= bbBlackKingCastle
		b.bbBlack ^= bbBlackKingCastle
		b.bbOccupied ^= bbBlackKingCastleTravel
	case c == Black && m.HasTag(QueenSideCastle): // black queen side castle
		b.bbRook ^= bbBlackQueenCastle
		b.bbBlack ^= bbBlackQueenCastle
		b.bbOccupied ^= bbBlackQueenCastleTravel
	default: // quiet
		b.bbOccupied ^= mbb
	}

	b.computeConvenienceBitboards()
}

func (b board) hasSufficientMaterial() bool {
	if (b.bbWhite&b.bbQueen | b.bbWhite&b.bbRook | b.bbWhite&b.bbPawn |
		b.bbBlack&b.bbQueen | b.bbBlack&b.bbRook | b.bbBlack&b.bbPawn) > 0 {
		return true
	}

	nWhites, nBlacks := b.bbWhite.ones(), b.bbBlack.ones()
	nWhiteBishops := (b.bbWhite & b.bbBishop).ones()
	nBlackBishops := (b.bbBlack & b.bbBishop).ones()

	// king versus king
	// king and bishop versus king
	// king and knight versus king
	if b.bbWhite == b.bbWhite&b.bbKing && b.bbBlack == b.bbBlack&b.bbKing ||
		b.bbWhite == b.bbWhite&b.bbKing && nBlacks == 2 && (nBlackBishops == 1 || (b.bbBlack&b.bbKnight).ones() == 1) ||
		b.bbBlack == b.bbBlack&b.bbKing && nWhites == 2 && (nWhiteBishops == 1 || (b.bbWhite&b.bbKnight).ones() == 1) {
		return false
	}

	// king and bishop versus king and bishop with the bishops on the same color
	if nWhites == 2 && nBlacks == 2 && nWhiteBishops == 1 && nBlackBishops == 1 &&
		((b.bbBishop&bbWhiteSquares).ones() == 2 ||
			(b.bbBishop&bbBlackSquares).ones() == 2) {
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
			if p := b.pieceAt(NewSquare(f, r)); p != NoPiece {
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

func (b board) getBitboard(p Piece) bitboard {
	switch p {
	case WhiteKing:
		return b.bbWhite & b.bbKing
	case WhiteQueen:
		return b.bbWhite & b.bbQueen
	case WhiteRook:
		return b.bbWhite & b.bbRook
	case WhiteBishop:
		return b.bbWhite & b.bbBishop
	case WhiteKnight:
		return b.bbWhite & b.bbKnight
	case WhitePawn:
		return b.bbWhite & b.bbPawn
	case BlackKing:
		return b.bbBlack & b.bbKing
	case BlackQueen:
		return b.bbBlack & b.bbQueen
	case BlackRook:
		return b.bbBlack & b.bbRook
	case BlackBishop:
		return b.bbBlack & b.bbBishop
	case BlackKnight:
		return b.bbBlack & b.bbKnight
	case BlackPawn:
		return b.bbBlack & b.bbPawn
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
	return (b.bbKing & b.getColor(c)).scanForward()
}

func (b *board) xorBitboard(pt PieceType, bb bitboard) {
	switch pt {
	case King:
		b.bbKing ^= bb
	case Queen:
		b.bbQueen ^= bb
	case Rook:
		b.bbRook ^= bb
	case Bishop:
		b.bbBishop ^= bb
	case Knight:
		b.bbKnight ^= bb
	case Pawn:
		b.bbPawn ^= bb
	}
}

func (b *board) xorColor(c Color, bb bitboard) {
	switch c {
	case White:
		b.bbWhite ^= bb
	case Black:
		b.bbBlack ^= bb
	}
}

func (b board) copyBoard() board {
	return board{
		bbKing:     b.bbKing,
		bbQueen:    b.bbQueen,
		bbRook:     b.bbRook,
		bbBishop:   b.bbBishop,
		bbKnight:   b.bbKnight,
		bbPawn:     b.bbPawn,
		bbWhite:    b.bbWhite,
		bbBlack:    b.bbBlack,
		bbOccupied: b.bbOccupied,
		bbPinned:   b.bbPinned,
		bbPinner:   b.bbPinner,
		bbCheck:    b.bbCheck,
	}
}
