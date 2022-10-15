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

// String returns a 64 character string of 1s and 0s starting with the most significant bit.
func (b bitboard) String() string {
	return fmt.Sprintf("%064b", b)
}
