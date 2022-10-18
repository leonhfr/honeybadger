package chess

import (
	"fmt"
	"math/bits"
)

type squareSet map[Square]struct{}

// bitboard is a board representation encoded in an unsigned 64-bit integer. The
// 64 board positions have A1 as the least significant bit and H8 as the most.
type bitboard uint64

func newBitboard(s squareSet) bitboard {
	var bb bitboard
	for sq := range s {
		bb |= 1 << sq
	}
	return bb
}

func (b bitboard) mapping() squareSet {
	s := squareSet{}
	for sq := A1; sq <= H8; sq++ {
		if b.occupied(sq) {
			s[sq] = struct{}{}
		}
	}
	return s
}

func (b bitboard) reverse() bitboard {
	return bitboard(bits.Reverse64(uint64(b)))
}

func (b bitboard) occupied(sq Square) bool {
	return (b & (1 << sq)) > 0
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

var (
	bbRanks = [8]bitboard{bbRank1, bbRank2, bbRank3, bbRank4, bbRank5, bbRank6, bbRank7, bbRank8}
	bbFiles = [8]bitboard{bbFileA, bbFileB, bbFileC, bbFileD, bbFileE, bbFileF, bbFileG, bbFileH}

	bbDiagonals     = [64]bitboard{}
	bbAntiDiagonals = [64]bitboard{}
	bbKingMoves     = [64]bitboard{}
	bbKnightMoves   = [64]bitboard{}
)

func init() {
	for sq := A1; sq <= H8; sq++ {
		bbDiagonals[sq] = initDiagonalBitboard(sq)
		bbAntiDiagonals[sq] = initAntiDiagonalBitboard(sq)
		bbKingMoves[sq] = initKingBitboard(sq)
		bbKnightMoves[sq] = initKnightBitboard(sq)
	}
}

func initKingBitboard(sq Square) bitboard {
	set := squareSet{}
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
	return newBitboard(set)
}

func initKnightBitboard(sq Square) bitboard {
	set := squareSet{}
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
	return newBitboard(set)
}

func initDiagonalBitboard(sq Square) bitboard {
	set := squareSet{}
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
	return newBitboard(set)
}

func initAntiDiagonalBitboard(sq Square) bitboard {
	set := squareSet{}
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
	return newBitboard(set)
}
