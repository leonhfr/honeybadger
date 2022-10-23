package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBitboard(t *testing.T) {
	for _, p := range pieces {
		t.Run(p.String(), func(t *testing.T) {
			assert.Equal(t, startingBoard.getBitboard(p), newBitboard(startingBoardMap[p]))
		})
	}
}

func TestBitboard_Mapping(t *testing.T) {
	for _, p := range pieces {
		t.Run(p.String(), func(t *testing.T) {
			assert.ElementsMatch(t, startingBoardMap[p], startingBoard.getBitboard(p).mapping())
		})
	}
}

func BenchmarkBitboard_Mapping(b *testing.B) {
	for _, p := range pieces {
		b.Run(p.String(), func(b *testing.B) {
			bb := startingBoard.getBitboard(p)
			for n := 0; n < b.N; n++ {
				bb.mapping()
			}
		})
	}
}

func TestBitboard_Reverse(t *testing.T) {
	tests := []struct {
		name string
		args bitboard
		want bitboard
	}{
		{"0", 0, 0},
		{"1", 1, 9223372036854775808},
		{"18446744073709551615", 18446744073709551615, 18446744073709551615},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.reverse())
		})
	}
}

func BenchmarkBitboard_Reverse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bitboard(9223372036854775807).reverse()
	}
}

func TestBitboard_Occupied(t *testing.T) {
	type args struct {
		bb bitboard
		sq Square
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{"A1", args{1, A1}, true},
		{"A1", args{1, A8}, false},
		{"A1", args{1, H1}, false},
		{"A1", args{1, H8}, false},
		{"A1A8H1H8", args{9295429630892703873, A1}, true},
		{"A1A8H1H8", args{9295429630892703873, A8}, true},
		{"A1A8H1H8", args{9295429630892703873, H1}, true},
		{"A1A8H1H8", args{9295429630892703873, H8}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.bb.occupied(tt.args.sq))
		})
	}
}
