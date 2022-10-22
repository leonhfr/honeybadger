package chess

import (
	"fmt"
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

func (b bitboard) mapping() []Square {
	s := []Square{}
	for sq := A1; sq <= H8; sq++ {
		if b.occupied(sq) {
			s = append(s, sq)
		}
	}
	return s
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
)

func init() {
	for sq := A1; sq <= H8; sq++ {
		bbRanks[sq] = initRankBitboard(sq)
		bbFiles[sq] = initFileBitboard(sq)
		bbDiagonals[sq] = initDiagonalBitboard(sq)
		bbAntiDiagonals[sq] = initAntiDiagonalBitboard(sq)
		bbKingMoves[sq] = initKingBitboard(sq)
		bbKnightMoves[sq] = initKnightBitboard(sq)
		bbWhitePawnPushes[sq] = initWhitePawnPushBitboard(sq)
		bbBlackPawnPushes[sq] = initBlackPawnPushBitboard(sq)
		bbWhitePawnCaptures[sq] = initWhitePawnCaptureBitboard(sq)
		bbBlackPawnCaptures[sq] = initBlackPawnCaptureBitboard(sq)
		bbDoubleSquares[sq] = initDoubleBitboard(sq)
		bbReverseDoubleSquares[sq] = initReverseDoubleBitboard(sq)
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

func squareSetToSlice(set map[Square]struct{}) []Square {
	s := []Square{}
	for sq := range set {
		s = append(s, sq)
	}
	return s
}
