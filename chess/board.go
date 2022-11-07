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
		bb := sq.bitboard()
		b.xorBitboard(p, bb)
		b.xorColor(p.Color(), bb)
		b.bbOccupied ^= bb
	}
	b.computeConvenienceBitboards()
	return b
}

func (b *board) computeConvenienceBitboards() {
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
		b.xorBitboard(p1, mbb)
	} else {
		// promotion
		b.xorBitboard(p1, s1bb)
		b.xorBitboard(promo, s2bb)
	}

	b.xorColor(c, mbb)

	switch enPassant := m.HasTag(EnPassant); {
	case m.HasTag(Capture) && !enPassant:
		// capture
		b.xorBitboard(p2, s2bb)
		b.xorColor(p2.Color(), s2bb)
		b.bbOccupied ^= s1bb
	case c == White && enPassant:
		// white en passant
		bb := (s2 - 8).bitboard()
		b.bbBlackPawn ^= bb
		b.bbBlack ^= bb
		b.bbOccupied ^= mbb ^ bb
	case c == Black && enPassant:
		// black en passant
		bb := (s2 + 8).bitboard()
		b.bbWhitePawn ^= bb
		b.bbWhite ^= bb
		b.bbOccupied ^= mbb ^ bb
	case c == White && m.HasTag(KingSideCastle):
		// white king side castle
		b.bbWhiteRook ^= bbWhiteKingCastle
		b.bbWhite ^= bbWhiteKingCastle
		b.bbOccupied ^= bbWhiteKingCastleTravel
	case c == White && m.HasTag(QueenSideCastle):
		// white queen side castle
		b.bbWhiteRook ^= bbWhiteQueenCastle
		b.bbWhite ^= bbWhiteQueenCastle
		b.bbOccupied ^= bbWhiteQueenCastleTravel
	case c == Black && m.HasTag(KingSideCastle):
		b.bbBlackRook ^= bbBlackKingCastle
		b.bbBlack ^= bbBlackKingCastle
		b.bbOccupied ^= bbBlackKingCastleTravel
	case c == Black && m.HasTag(QueenSideCastle):
		// black queen side castle
		b.bbBlackRook ^= bbBlackQueenCastle
		b.bbBlack ^= bbBlackQueenCastle
		b.bbOccupied ^= bbBlackQueenCastleTravel
	default:
		// quiet
		b.bbOccupied ^= mbb
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

func (b *board) xorBitboard(p Piece, bb bitboard) {
	switch p {
	case WhiteKing:
		b.bbWhiteKing ^= bb
	case WhiteQueen:
		b.bbWhiteQueen ^= bb
	case WhiteRook:
		b.bbWhiteRook ^= bb
	case WhiteBishop:
		b.bbWhiteBishop ^= bb
	case WhiteKnight:
		b.bbWhiteKnight ^= bb
	case WhitePawn:
		b.bbWhitePawn ^= bb
	case BlackKing:
		b.bbBlackKing ^= bb
	case BlackQueen:
		b.bbBlackQueen ^= bb
	case BlackRook:
		b.bbBlackRook ^= bb
	case BlackBishop:
		b.bbBlackBishop ^= bb
	case BlackKnight:
		b.bbBlackKnight ^= bb
	case BlackPawn:
		b.bbBlackPawn ^= bb
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
