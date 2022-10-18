package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiagonalMoves(t *testing.T) {
	occupied := squareSet{
		F6: struct{}{},
		B6: struct{}{},
	}
	want := squareSet{
		B6: struct{}{},
		F6: struct{}{},
		C5: struct{}{},
		E5: struct{}{},
		C3: struct{}{},
		E3: struct{}{},
		B2: struct{}{},
		F2: struct{}{},
		A1: struct{}{},
		G1: struct{}{},
	}
	assert.Equal(t, want, diagonalMoves(D4, newBitboard(occupied)).mapping())
}

func TestHVMoves(t *testing.T) {
	occupied := squareSet{
		D3: struct{}{},
		F5: struct{}{},
	}
	want := squareSet{
		D8: struct{}{},
		D7: struct{}{},
		D6: struct{}{},
		A5: struct{}{},
		B5: struct{}{},
		C5: struct{}{},
		E5: struct{}{},
		F5: struct{}{},
		D4: struct{}{},
		D3: struct{}{},
	}
	assert.Equal(t, want, hvMoves(D5, newBitboard(occupied)).mapping())
}
