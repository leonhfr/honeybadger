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
	if b == 0 {
		return nil
	}
	squares := make([]Square, 0, 8)
	for b > 0 {
		squares = append(squares, b.scanForward())
		b = b.resetLSB()
	}
	return squares
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
