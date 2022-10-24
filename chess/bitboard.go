package chess

import (
	"fmt"
	"math"
	"math/bits"
)

// bitboard is a board representation encoded in an unsigned 64-bit integer. The
// 64 board positions have A1 as the least significant bit and H8 as the most.
type bitboard uint64

func newBitboard(s []Square) bitboard {
	var bb bitboard
	for _, sq := range s {
		bb |= sq.bitboard()
	}
	return bb
}

var deBruijnMap = [64]Square{
	A1, H6, B1, A8, A7, D4, C1, E8,
	B8, B7, B6, F5, E4, A3, D1, F8,
	G7, C8, D5, E7, C7, C6, F3, E6,
	G5, A5, F4, H3, B3, D2, E1, G8,
	G6, H7, C4, D8, A6, E5, H2, F7,
	C5, D7, E3, D6, H4, G3, C2, F6,
	B4, H5, G2, B5, D3, G4, B2, A4,
	F2, C3, A2, E2, H1, G1, F1, H8,
}

const deBruijn = 0x03f79d71b4cb0a89

func (b bitboard) mapping() []Square {
	if b == 0 {
		return nil
	}
	var squares []Square
	for b > 0 {
		squares = append(squares, b.scanForward())
		b = b.resetLSB()
	}
	return squares
}

// bitboard can't be 0
//
// uses de Bruijn forward scanning
func (b bitboard) scanForward() Square {
	i := ((b ^ (b - 1)) * deBruijn) >> 58
	return deBruijnMap[i]
}

// resets lowest significant bit
func (b bitboard) resetLSB() bitboard {
	return b & (b - 1)
}

func (b bitboard) reverse() bitboard {
	return bitboard(bits.Reverse64(uint64(b)))
}

func (b bitboard) occupied(sq Square) bool {
	return (b & sq.bitboard()) > 0
}

func (b bitboard) ones() int {
	return bits.OnesCount64(uint64(b))
}

// String returns a 64 character string of 1s and 0s starting with the most significant bit.
func (b bitboard) String() string {
	return fmt.Sprintf("%064b", b)
}

const (
	bbRank1 bitboard = (1<<A1 + 1<<B1 + 1<<C1 + 1<<D1 + 1<<E1 + 1<<F1 + 1<<G1 + 1<<H1) << (8 * iota)
	bbRank2
	bbRank3
	bbRank4
	bbRank5
	bbRank6
	bbRank7
	bbRank8
)

const (
	bbFileA bitboard = (1<<A1 + 1<<A2 + 1<<A3 + 1<<A4 + 1<<A5 + 1<<A6 + 1<<A7 + 1<<A8) << iota
	bbFileB
	bbFileC
	bbFileD
	bbFileE
	bbFileF
	bbFileG
	bbFileH
)

const (
	bbNotRank1 = ^bbRank1
	bbNotRank8 = ^bbRank8
	bbNotFileA = ^bbFileA
	bbNotFileH = ^bbFileH
)

const (
	bbWhiteSquares bitboard = 1<<B1 + 1<<D1 + 1<<F1 + 1<<H1 +
		1<<B3 + 1<<D3 + 1<<F3 + 1<<H3 +
		1<<B5 + 1<<D5 + 1<<F5 + 1<<H5 +
		1<<B7 + 1<<D7 + 1<<F7 + 1<<H7 +
		1<<A2 + 1<<C2 + 1<<E2 + 1<<G2 +
		1<<A4 + 1<<C4 + 1<<E4 + 1<<G4 +
		1<<A6 + 1<<C6 + 1<<E6 + 1<<G6 +
		1<<A8 + 1<<C8 + 1<<E8 + 1<<G8
	bbBlackSquares = ^bbWhiteSquares
)

var (
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
	return (sq.bitboard() & bbNotRank8) << 8
}

func initBlackPawnPushBitboard(sq Square) bitboard {
	return (sq.bitboard() & bbNotRank1) >> 8
}

func initWhitePawnCaptureBitboard(sq Square) bitboard {
	captureR := (sq.bitboard() & bbNotFileH & bbNotRank8) << 9
	captureL := (sq.bitboard() & bbNotFileA & bbNotRank8) << 7
	return captureR | captureL
}

func initBlackPawnCaptureBitboard(sq Square) bitboard {
	captureR := (sq.bitboard() & bbNotFileH & bbNotRank1) >> 7
	captureL := (sq.bitboard() & bbNotFileA & bbNotRank1) >> 9
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
