package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBitboard(t *testing.T) {
	tests := []struct {
		name string
		args []Square
		want string
	}{
		{
			"A1",
			[]Square{A1},
			"0000000000000000000000000000000000000000000000000000000000000001",
		},
		{
			"A1A8H1H8",
			[]Square{A1, A8, H1, H8},
			"1000000100000000000000000000000000000000000000000000000010000001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, newBitboard(tt.args).String())
		})
	}
}

func TestBitboard_Mapping(t *testing.T) {
	tests := []struct {
		name string
		args bitboard
		want []Square
	}{
		{
			"A1",
			1,
			[]Square{A1},
		},
		{
			"A1A8H1H8",
			bitboard(9295429630892703873),
			[]Square{A1, A8, H1, H8},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.want, tt.args.mapping())
		})
	}
}

func BenchmarkBitboard_Mapping(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bitboard(9223372036854775807).mapping()
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
