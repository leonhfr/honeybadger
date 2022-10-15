package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBitboard(t *testing.T) {
	tests := []struct {
		name string
		args squareMap
		want string
	}{
		{
			"A1",
			squareMap{A1: struct{}{}},
			"0000000000000000000000000000000000000000000000000000000000000001",
		},
		{
			"A1A8H1H8",
			squareMap{A1: struct{}{}, A8: struct{}{}, H1: struct{}{}, H8: struct{}{}},
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
		want squareMap
	}{
		{
			"A1",
			1,
			squareMap{A1: struct{}{}},
		},
		{
			"A1A8H1H8",
			bitboard(9295429630892703873),
			squareMap{A1: struct{}{}, A8: struct{}{}, H1: struct{}{}, H8: struct{}{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.mapping())
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
