package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMoveBitboard(t *testing.T) {
	fen := "k7/8/8/8/8/8/5P2/KQRBN3 w - - 0 1"
	pos := unsafeFEN(fen)

	type args struct {
		sq Square
		pt PieceType
	}
	tests := []struct {
		args args
		want squareSet
	}{
		{args{A1, King}, squareSet{
			A2: struct{}{}, B2: struct{}{},
			B1: struct{}{}, // will be removed
		}},
		{args{B1, Queen}, squareSet{
			A2: struct{}{}, B2: struct{}{},
			B3: struct{}{}, B4: struct{}{},
			B5: struct{}{}, B6: struct{}{},
			B7: struct{}{}, B8: struct{}{},
			C2: struct{}{}, D3: struct{}{},
			E4: struct{}{}, F5: struct{}{},
			G6: struct{}{}, H7: struct{}{},
			A1: struct{}{}, C1: struct{}{}, // will be removed
		}},
		{args{C1, Rook}, squareSet{
			C2: struct{}{}, C3: struct{}{},
			C4: struct{}{}, C5: struct{}{},
			C6: struct{}{}, C7: struct{}{},
			C8: struct{}{},
			B1: struct{}{}, D1: struct{}{}, // will be removed
		}},
		{args{D1, Bishop}, squareSet{
			A4: struct{}{}, B3: struct{}{},
			C2: struct{}{}, E2: struct{}{},
			F3: struct{}{}, G4: struct{}{},
			H5: struct{}{},
		}},
		{args{E1, Knight}, squareSet{
			C2: struct{}{}, D3: struct{}{},
			F3: struct{}{}, G2: struct{}{},
		}},
		{args{F2, Pawn}, squareSet{
			F3: struct{}{}, F4: struct{}{},
		}},
		{args{F2, NoPieceType}, squareSet{}},
	}

	for _, tt := range tests {
		t.Run(tt.args.pt.String(), func(t *testing.T) {
			got := moveBitboard(tt.args.sq, pos, tt.args.pt)
			assert.Equal(t, tt.want, got.mapping())
		})
	}
}

func TestPawnBitboards(t *testing.T) {
	fenWhite := "k7/p7/1p6/2N5/2n2pP1/1P6/P7/K7 w - - 0 1"
	fenBlack := "k7/p7/1p6/2N5/2n2pP1/1P6/P7/K7 b - g3 0 1"
	posWhite, posBlack := unsafeFEN(fenWhite), unsafeFEN(fenBlack)
	tests := []struct {
		sq  Square
		set squareSet
		pos *Position
	}{
		{A2, squareSet{A3: struct{}{}, A4: struct{}{}}, posWhite},
		{B3, squareSet{B4: struct{}{}, C4: struct{}{}}, posWhite},
		{A7, squareSet{A5: struct{}{}, A6: struct{}{}}, posBlack},
		{B6, squareSet{B5: struct{}{}, C5: struct{}{}}, posBlack},
		{F4, squareSet{F3: struct{}{}, G3: struct{}{}}, posBlack},
	}

	for _, tt := range tests {
		t.Run(tt.sq.String(), func(t *testing.T) {
			assert.Equal(t, tt.set, pawnBitboard(tt.sq, tt.pos).mapping())
		})
	}
}

func TestDiagonalBitboard(t *testing.T) {
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
	assert.Equal(t, want, diagonalBitboard(D4, newBitboard(occupied)).mapping())
}

func TestHVBitboard(t *testing.T) {
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
	assert.Equal(t, want, hvBitboard(D5, newBitboard(occupied)).mapping())
}
