package chess

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCastlingMoves(t *testing.T) {
	tests := []struct {
		args string
		want []string
	}{
		{"r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1", []string{"e1c1", "e1g1"}},
		{"r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1", []string{"e8c8", "e8g8"}},
	}

	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			var moves []string
			for _, m := range castlingMoves(unsafeFEN(tt.args)) {
				moves = append(moves, m.String())
			}
			sort.Strings(moves)
			assert.Equal(t, tt.want, moves)
		})
	}
}

func TestPseudoMoves(t *testing.T) {
	tests := []struct {
		args string
		want []string
	}{
		{
			"1k2q3/8/8/8/8/8/4R3/4K3 w - - 0 1",
			[]string{
				"e1d1", "e1d2", "e1f1", "e1f2", "e2a2",
				"e2b2", "e2c2", "e2d2", "e2e3", "e2e4",
				"e2e5", "e2e6", "e2e7", "e2e8", "e2f2",
				"e2g2", "e2h2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			var moves []string
			for _, m := range pseudoMoves(unsafeFEN(tt.args)) {
				moves = append(moves, m.String())
			}
			sort.Strings(moves)
			assert.Equal(t, tt.want, moves)
		})
	}
}

func TestIsAttacked(t *testing.T) {
	fen := "k6q/8/8/8/8/8/8/K7 w - - 0 1"
	pos := unsafeFEN(fen)
	assert.True(t, isAttacked(pos.board.sqWhiteKing, pos))
}

func TestIsAttackedByCount(t *testing.T) {
	fen := "K2r3q/8/8/2p5/r2Q4/2k2n2/4n3/6b1 w - - 0 1"
	pos := unsafeFEN(fen)
	sq := D4

	tests := []struct {
		args PieceType
		want int
	}{
		{King, 1},
		{Queen, 1},
		{Rook, 2},
		{Bishop, 1},
		{Knight, 2},
		{Pawn, 1},
	}

	for _, tt := range tests {
		t.Run(tt.args.String(), func(t *testing.T) {
			assert.Equal(t, tt.want, isAttackedByCount(sq, pos, tt.args))
		})
	}
}

func TestIsAttackedByPawnCount(t *testing.T) {
	type args struct {
		sq  Square
		fen string
	}

	tests := []struct {
		args args
		want int
	}{
		{args{A2, "k7/p1p5/1P1P3P/8/8/1p1p3p/P1P5/K7 w - - 0 1"}, 1},
		{args{C2, "k7/p1p5/1P1P3P/8/8/1p1p3p/P1P5/K7 w - - 0 1"}, 2},
		{args{A7, "k7/p1p5/1P1P3P/8/8/1p1p3p/P1P5/K7 b - - 0 1"}, 1},
		{args{C7, "k7/p1p5/1P1P3P/8/8/1p1p3p/P1P5/K7 b - - 0 1"}, 2},
		{args{G4, "k7/p1p5/1P1P3P/5Pp1/5pP1/1p1p3p/P1P5/K7 w - g3 0 1"}, 1},
		{args{G5, "k7/p1p5/1P1P3P/5Pp1/5pP1/1p1p3p/P1P5/K7 b - g6 0 1"}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.args.sq.String(), func(t *testing.T) {
			got := isAttackedByPawnCount(tt.args.sq, unsafeFEN(tt.args.fen))
			assert.Equal(t, tt.want, got)
		})
	}
}

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
