package chess

import "math"

type castleCheck struct {
	bbPawn   bitboard
	bbKnight bitboard
	bbKing   bitboard
	squares  [3]Square
}

var (
	castleChecks = [4]castleCheck{}

	bbRanks                = [64]bitboard{}
	bbFiles                = [64]bitboard{}
	bbDiagonals            = [64]bitboard{}
	bbAntiDiagonals        = [64]bitboard{}
	bbKingMoves            = [64]bitboard{}
	bbKnightMoves          = [64]bitboard{}
	bbWhitePawnPushes      = [64]bitboard{}
	bbBlackPawnPushes      = [64]bitboard{}
	bbWhitePawnCaptures    = [64]bitboard{}
	bbBlackPawnCaptures    = [64]bitboard{}
	bbDoubleSquares        = [64]bitboard{}
	bbReverseDoubleSquares = [64]bitboard{}
	bbInBetween            = [64][64]bitboard{}
)

func init() {
	for s1 := A1; s1 <= H8; s1++ {
		bbRanks[s1] = initRankBitboard(s1)
		bbFiles[s1] = initFileBitboard(s1)
		bbDiagonals[s1] = initDiagonalBitboard(s1)
		bbAntiDiagonals[s1] = initAntiDiagonalBitboard(s1)
		bbKingMoves[s1] = initKingBitboard(s1)
		bbKnightMoves[s1] = initKnightBitboard(s1)
		bbWhitePawnPushes[s1] = initWhitePawnPushBitboard(s1)
		bbBlackPawnPushes[s1] = initBlackPawnPushBitboard(s1)
		bbWhitePawnCaptures[s1] = initWhitePawnCaptureBitboard(s1)
		bbBlackPawnCaptures[s1] = initBlackPawnCaptureBitboard(s1)
		bbDoubleSquares[s1] = initDoubleBitboard(s1)
		bbReverseDoubleSquares[s1] = initReverseDoubleBitboard(s1)

		for s2 := A1; s2 <= H8; s2++ {
			bbInBetween[s1][s2] = initInBetweenBitboard(s1, s2)
		}
	}

	initCastleChecks()
}

func initCastleChecks() {
	castles := [4]struct {
		color   Color
		squares [3]Square
	}{
		{White, [3]Square{E1, F1, G1}},
		{White, [3]Square{C1, D1, E1}},
		{Black, [3]Square{E8, F8, G8}},
		{Black, [3]Square{C8, D8, E8}},
	}

	for i, castle := range castles {
		var bbPawn, bbKnight, bbKing bitboard
		for _, sq := range castle.squares {
			if castle.color == White {
				bbPawn |= bbWhitePawnCaptures[sq]
			} else {
				bbPawn |= bbBlackPawnCaptures[sq]
			}
			bbKnight |= bbKnightMoves[sq]
			bbKing |= bbKingMoves[sq]
		}

		castleChecks[i] = castleCheck{bbPawn, bbKnight, bbKing, castle.squares}
	}
}

func initRankBitboard(sq Square) bitboard {
	bbRanks := [8]bitboard{bbRank1, bbRank2, bbRank3, bbRank4, bbRank5, bbRank6, bbRank7, bbRank8}
	return bbRanks[sq.Rank()/8]
}

func initFileBitboard(sq Square) bitboard {
	bbFiles := [8]bitboard{bbFileA, bbFileB, bbFileC, bbFileD, bbFileE, bbFileF, bbFileG, bbFileH}
	return bbFiles[sq.File()]
}

func initDiagonalBitboard(sq Square) bitboard {
	set := map[Square]struct{}{}
	for i := 0; i < 8; i++ {
		upLeft := sq + 8*Square(i) - Square(i)
		set[upLeft] = struct{}{}
		if upLeft.File() == FileA || upLeft.Rank() == Rank8 {
			break
		}
	}
	for i := 0; i < 8; i++ {
		downRight := sq - 8*Square(i) + Square(i)
		set[downRight] = struct{}{}
		if downRight.File() == FileH || downRight.Rank() == Rank1 {
			break
		}
	}
	return newBitboard(squareSetToSlice(set))
}

func initAntiDiagonalBitboard(sq Square) bitboard {
	set := map[Square]struct{}{}
	for i := 0; i < 8; i++ {
		upRight := sq + 8*Square(i) + Square(i)
		set[upRight] = struct{}{}
		if upRight.File() == FileH || upRight.Rank() == Rank8 {
			break
		}
	}
	for i := 0; i < 8; i++ {
		downLeft := sq - 8*Square(i) - Square(i)
		set[downLeft] = struct{}{}
		if downLeft.File() == FileA || downLeft.Rank() == Rank1 {
			break
		}
	}
	return newBitboard(squareSetToSlice(set))
}

func initKingBitboard(sq Square) bitboard {
	set := map[Square]struct{}{}
	for dest, ok := range map[Square]bool{
		sq + 8 - 1: sq.Rank() <= Rank7 && sq.File() >= FileB,
		sq + 8:     sq.Rank() <= Rank7,
		sq + 8 + 1: sq.Rank() <= Rank7 && sq.File() <= FileG,
		sq + 1:     sq.File() <= FileG,
		sq - 8 + 1: sq.Rank() >= Rank2 && sq.File() <= FileG,
		sq - 8:     sq.Rank() >= Rank2,
		sq - 8 - 1: sq.Rank() >= Rank2 && sq.File() >= FileB,
		sq - 1:     sq.File() >= FileB,
	} {
		if ok {
			set[dest] = struct{}{}
		}
	}
	return newBitboard(squareSetToSlice(set))
}

func initKnightBitboard(sq Square) bitboard {
	set := map[Square]struct{}{}
	for dest, ok := range map[Square]bool{
		sq + 8 - 2:  sq.Rank() <= Rank7 && sq.File() >= FileC,
		sq + 16 - 1: sq.Rank() <= Rank6 && sq.File() >= FileB,
		sq + 16 + 1: sq.Rank() <= Rank6 && sq.File() <= FileG,
		sq + 8 + 2:  sq.Rank() <= Rank7 && sq.File() <= FileF,
		sq - 8 + 2:  sq.Rank() >= Rank2 && sq.File() <= FileF,
		sq - 16 + 1: sq.Rank() >= Rank3 && sq.File() <= FileG,
		sq - 16 - 1: sq.Rank() >= Rank3 && sq.File() >= FileB,
		sq - 8 - 2:  sq.Rank() >= Rank2 && sq.File() >= FileC,
	} {
		if ok {
			set[dest] = struct{}{}
		}
	}
	return newBitboard(squareSetToSlice(set))
}

func initWhitePawnPushBitboard(sq Square) bitboard {
	return (sq.bitboard() & ^bbRank8) << 8
}

func initBlackPawnPushBitboard(sq Square) bitboard {
	return (sq.bitboard() & ^bbRank1) >> 8
}

func initWhitePawnCaptureBitboard(sq Square) bitboard {
	captureR := (sq.bitboard() & ^bbFileH & ^bbRank8) << 9
	captureL := (sq.bitboard() & ^bbFileA & ^bbRank8) << 7
	return captureR | captureL
}

func initBlackPawnCaptureBitboard(sq Square) bitboard {
	captureR := (sq.bitboard() & ^bbFileH & ^bbRank1) >> 7
	captureL := (sq.bitboard() & ^bbFileA & ^bbRank1) >> 9
	return captureR | captureL
}

func initDoubleBitboard(sq Square) bitboard {
	return 2 * sq.bitboard()
}

func initReverseDoubleBitboard(sq Square) bitboard {
	return 2 * sq.bitboard().reverse()
}

func initInBetweenBitboard(s1, s2 Square) bitboard {
	var m1, a2a7, b2g7, h1b7 bitboard = math.MaxUint64,
		0x0001010101010100, 0x0040201008040200, 0x0002040810204080

	between := (m1 >> (64 - s2)) ^ (m1 >> (64 - s1))
	file := bitboard((s2 & 7) - (s1 & 7))
	rank := bitboard(((s2 | 7) - s1) >> 3)
	line := ((file & 7) - 1) & a2a7
	line += 2 * (((rank & 7) - 1) >> 58)
	line += (((rank - file) & 15) - 1) & b2g7
	line += (((rank + file) & 15) - 1) & h1b7
	line *= between & -between
	return line & between
}

func squareSetToSlice(set map[Square]struct{}) []Square {
	s := []Square{}
	for sq := range set {
		s = append(s, sq)
	}
	return s
}
