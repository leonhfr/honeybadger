package chess

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

type perfTest struct {
	fen   string
	nodes []int
}

var perfResults = []perfTest{
	{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		[]int{
			20, 400, 8902, 197281,
			// 4865609, 119060324, 3195901860, 84998978956, 2439530234167, 69352859712417
		},
	},
	{
		"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
		[]int{
			48, 2039, 97862,
			// 4085603, 193690690
		},
	},
	{
		"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
		[]int{
			14, 191, 2812, 43238, 674624,
			// 11030083, 178633661
		},
	},
	{
		"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
		[]int{
			6, 264, 9467, 422333,
			// 15833292, 706045033
		},
	},
	{
		"r2q1rk1/pP1p2pp/Q4n2/bbp1p3/Np6/1B3NBn/pPPP1PPP/R3K2R b KQ - 0 1",
		[]int{
			6, 264, 9467, 422333,
			// 15833292, 706045033
		},
	},
	{
		"rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
		[]int{
			44, 1486, 62379, 2103487,
			// 89941194,
		},
	},
	{
		"r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10",
		[]int{
			46, 2079, 89890,
			// 3894594, 164075551, 6923051137, 287188994746, 11923589843526, 490154852788714
		},
	},
}

func TestPerfResults(t *testing.T) {
	for _, tt := range perfResults {
		t.Run(tt.fen, func(t *testing.T) {
			for depth := 0; depth < len(tt.nodes); depth++ {
				want := tt.nodes[depth]
				got := perft(unsafeFEN(tt.fen), depth)
				assert.Equal(t, want, got)
			}
		})
	}
}

func perft(pos *Position, depth int) int {
	if depth == 0 {
		return len(legalMoves(pos))
	}

	var count int
	for _, m := range legalMoves(pos) {
		count += perft(pos.Move(m), depth-1)
	}
	return count
}

func BenchmarkLegalMoves(b *testing.B) {
	pos := StartingPosition()
	for n := 0; n < b.N; n++ {
		legalMoves(pos)
	}
}

func BenchmarkPseudoMoves(b *testing.B) {
	pos := StartingPosition()
	for n := 0; n < b.N; n++ {
		pseudoMoves(pos)
	}
}

func TestCastlingMoves(t *testing.T) {
	tests := []struct {
		args string
		want []string
	}{
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", nil},
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

func benchmarkAttackByCount(pt PieceType, b *testing.B) {
	pos := unsafeFEN("K2r3q/8/8/2p5/r2Q4/2k2n2/4n3/6b1 w - - 0 1")
	sq := D4
	for n := 0; n < b.N; n++ {
		isAttackedByCount(sq, pos, pt)
	}
}

func BenchmarkAttackByCountKing(b *testing.B)   { benchmarkAttackByCount(King, b) }
func BenchmarkAttackByCountQueen(b *testing.B)  { benchmarkAttackByCount(Queen, b) }
func BenchmarkAttackByCountRook(b *testing.B)   { benchmarkAttackByCount(Rook, b) }
func BenchmarkAttackByCountBishop(b *testing.B) { benchmarkAttackByCount(Bishop, b) }
func BenchmarkAttackByCountKnight(b *testing.B) { benchmarkAttackByCount(Knight, b) }
func BenchmarkAttackByCountPawn(b *testing.B)   { benchmarkAttackByCount(Pawn, b) }

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
		want []Square
	}{
		{args{A1, King}, []Square{
			A2, B2,
			B1, // will be removed
		}},
		{args{B1, Queen}, []Square{
			A2, B2, B3, B4, B5,
			B6, B7, B8, C2, D3,
			E4, F5, G6, H7,
			A1, C1, // will be removed
		}},
		{args{C1, Rook}, []Square{
			C2, C3, C4, C5, C6,
			C7, C8,
			B1, D1, // will be removed
		}},
		{args{D1, Bishop}, []Square{
			A4, B3, C2, E2, F3,
			G4, H5,
		}},
		{args{E1, Knight}, []Square{C2, D3, F3, G2}},
		{args{F2, Pawn}, []Square{F3, F4}},
		{args{F2, NoPieceType}, []Square{}},
	}

	for _, tt := range tests {
		t.Run(tt.args.pt.String(), func(t *testing.T) {
			got := moveBitboard(tt.args.sq, pos, tt.args.pt)
			assert.ElementsMatch(t, tt.want, got.mapping())
		})
	}
}

func benchmarkMoveBitboard(sq Square, pt PieceType, b *testing.B) {
	pos := unsafeFEN("k7/8/8/8/8/8/5P2/KQRBN3 w - - 0 1")
	for n := 0; n < b.N; n++ {
		moveBitboard(sq, pos, pt)
	}
}

func BenchmarkMoveBitboardKing(b *testing.B)   { benchmarkMoveBitboard(A1, King, b) }
func BenchmarkMoveBitboardQueen(b *testing.B)  { benchmarkMoveBitboard(B1, Queen, b) }
func BenchmarkMoveBitboardRook(b *testing.B)   { benchmarkMoveBitboard(C1, Rook, b) }
func BenchmarkMoveBitboardBishop(b *testing.B) { benchmarkMoveBitboard(D1, Bishop, b) }
func BenchmarkMoveBitboardKnight(b *testing.B) { benchmarkMoveBitboard(E1, Knight, b) }
func BenchmarkMoveBitboardPawn(b *testing.B)   { benchmarkMoveBitboard(F2, Pawn, b) }

func TestPawnBitboards(t *testing.T) {
	fenWhite := "k7/p7/1p6/2N5/2n2pP1/1P6/P7/K7 w - - 0 1"
	fenBlack := "k7/p7/1p6/2N5/2n2pP1/1P6/P7/K7 b - g3 0 1"
	posWhite, posBlack := unsafeFEN(fenWhite), unsafeFEN(fenBlack)
	tests := []struct {
		sq  Square
		set []Square
		pos *Position
	}{
		{A2, []Square{A3, A4}, posWhite},
		{B3, []Square{B4, C4}, posWhite},
		{A7, []Square{A5, A6}, posBlack},
		{B6, []Square{B5, C5}, posBlack},
		{F4, []Square{F3, G3}, posBlack},
	}

	for _, tt := range tests {
		t.Run(tt.sq.String(), func(t *testing.T) {
			assert.ElementsMatch(t, tt.set, pawnBitboard(tt.sq, tt.pos).mapping())
		})
	}
}

func TestDiagonalBitboard(t *testing.T) {
	occupied := []Square{F6, B6}
	want := []Square{
		B6, F6, C5, E5, C3,
		E3, B2, F2, A1, G1,
	}
	assert.ElementsMatch(t, want, diagonalBitboard(D4, newBitboard(occupied)).mapping())
}

func TestHVBitboard(t *testing.T) {
	occupied := []Square{D3, F5}
	want := []Square{
		D8, D7, D6, A5, B5,
		C5, E5, F5, D4, D3,
	}
	assert.ElementsMatch(t, want, hvBitboard(D5, newBitboard(occupied)).mapping())
}
