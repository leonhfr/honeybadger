package chess

import (
	"fmt"
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
			20, 400, 8902, 197281, 4865609,
			// 119060324, 3195901860, 84998978956, 2439530234167, 69352859712417
		},
	},
	{
		"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
		[]int{
			48, 2039, 97862, 4085603,
			//  193690690
		},
	},
	{
		"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
		[]int{
			14, 191, 2812, 43238, 674624, 11030083,
			//  178633661
		},
	},
	{
		"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
		[]int{
			6, 264, 9467, 422333, 15833292,
			// 706045033
		},
	},
	{
		"r2q1rk1/pP1p2pp/Q4n2/bbp1p3/Np6/1B3NBn/pPPP1PPP/R3K2R b KQ - 0 1",
		[]int{
			6, 264, 9467, 422333, 15833292,
			// 706045033
		},
	},
	{
		"rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
		[]int{
			44, 1486, 62379, 2103487, 89941194,
		},
	},
	{
		"r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10",
		[]int{
			46, 2079, 89890, 3894594,
			//  164075551, 6923051137, 287188994746, 11923589843526, 490154852788714
		},
	},
}

func TestPerfResults(t *testing.T) {
	for _, tt := range perfResults {
		for depth := 0; depth < len(tt.nodes); depth++ {
			want := tt.nodes[depth]

			if !testing.Short() || want < 2<<22 {
				t.Run(fmt.Sprintf("%s depth %d", tt.fen, depth), func(t *testing.T) {
					got := perft(unsafeFEN(tt.fen), depth)
					assert.Equal(t, want, got)
				})
			}
		}
	}
}

func perft(pos *Position, depth int) int {
	if depth == 0 {
		return len(legalMoves(pos))
	}

	var count int
	for _, m := range pseudoMoves(pos) {
		if meta, ok := pos.MakeMove(m); ok {
			count += perft(pos, depth-1)
			pos.UnmakeMove(m, meta)
		}
	}
	return count
}

func legalMoves(pos *Position) []Move {
	var moves []Move
	for _, m := range pseudoMoves(pos) {
		if meta, ok := pos.MakeMove(m); ok {
			moves = append(moves, m)
			pos.UnmakeMove(m, meta)
		}
	}
	return moves
}

func BenchmarkPseudoMoves(b *testing.B) {
	for _, bb := range testPositions {
		pos := unsafeFEN(bb.preFEN)
		b.Run(bb.preFEN, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				pseudoMoves(pos)
			}
		})
	}
}

func TestCheckAttackAndInterposingMoves(t *testing.T) {
	tests := []struct {
		args string
		want []Move
	}{
		{
			"8/5pB1/3bn3/2r5/1q6/8/1k4K1/8 b - - 0 1",
			[]Move{
				newMove(BlackQueen, NoPiece, B4, C3, NoSquare, NoPiece),
				newMove(BlackQueen, NoPiece, B4, D4, NoSquare, NoPiece),
				newMove(BlackRook, NoPiece, C5, C3, NoSquare, NoPiece),
				newMove(BlackRook, NoPiece, C5, E5, NoSquare, NoPiece),
				newMove(BlackBishop, NoPiece, D6, E5, NoSquare, NoPiece),
				newMove(BlackKnight, NoPiece, E6, D4, NoSquare, NoPiece),
				newMove(BlackKnight, WhiteBishop, E6, G7, NoSquare, NoPiece),
				newMove(BlackPawn, NoPiece, F7, F6, NoSquare, NoPiece),
			},
		},
		{
			"8/8/8/8/4n3/2B5/1k4K1/8 b - - 0 1",
			[]Move{
				newMove(BlackKnight, WhiteBishop, E4, C3, NoSquare, NoPiece),
			},
		},
		{ // en passant defense
			"8/8/3p4/1Pp2r2/1K1R1p1k/8/4P1P1/8 w - c6 0 1",
			[]Move{
				newMove(WhitePawn, NoPiece, B5, C6, C6, NoPiece),
			},
		},
		{ // promotion
			"8/2p5/3p4/KP5r/4P2k/8/6p1/7R b - - 0 1",
			[]Move{
				newMove(BlackPawn, WhiteRook, G2, H1, NoSquare, BlackQueen),
				newMove(BlackPawn, WhiteRook, G2, H1, NoSquare, BlackRook),
				newMove(BlackPawn, WhiteRook, G2, H1, NoSquare, BlackBishop),
				newMove(BlackPawn, WhiteRook, G2, H1, NoSquare, BlackKnight),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			pos := unsafeFEN(tt.args)
			moves := checkAttackAndInterposingMoves(pos)
			assert.ElementsMatch(t, tt.want, moves)
		})
	}
}

func BenchmarkCheckAttackAndInterposingMoves(b *testing.B) {
	fens := []string{
		"8/5pB1/3bn3/2r5/1q6/8/1k4K1/8 b - - 0 1",
		"8/8/8/8/4n3/2B5/1k4K1/8 b - - 0 1",
		"8/8/3p4/1Pp2r2/1K1R1p1k/8/4P1P1/8 w - c6 0 1",
		"8/2p5/3p4/KP5r/4P2k/8/6p1/7R b - - 0 1",
	}

	for _, fen := range fens {
		pos := unsafeFEN(fen)
		b.Run(fen, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				checkAttackAndInterposingMoves(pos)
			}
		})
	}
}

func TestCheckFlightMoves(t *testing.T) {
	tests := []struct {
		args string
		want []Move
	}{
		{
			"2r2q1k/5pp1/4p1N1/8/1bp5/5P1R/6P1/2R4K b - - 0 1",
			[]Move{
				newMove(BlackKing, NoPiece, H8, G8, NoSquare, NoPiece),
			},
		},
		{
			"rnbk1b1r/pp3ppp/2p5/4q1B1/4n3/8/PPP2PPP/2KR1BNR b - - 0 1",
			[]Move{
				newMove(BlackKing, NoPiece, D8, E8, NoSquare, NoPiece),
				newMove(BlackKing, NoPiece, D8, C7, NoSquare, NoPiece),
			},
		},
		{
			"8/1Kr5/8/8/8/8/6k1/8 w - - 0 1",
			[]Move{
				newMove(WhiteKing, NoPiece, B7, A6, NoSquare, NoPiece),
				newMove(WhiteKing, NoPiece, B7, B6, NoSquare, NoPiece),
				newMove(WhiteKing, BlackRook, B7, C7, NoSquare, NoPiece),
				newMove(WhiteKing, NoPiece, B7, A8, NoSquare, NoPiece),
				newMove(WhiteKing, NoPiece, B7, B8, NoSquare, NoPiece),
			},
		},
		{ // en passant defense
			"8/8/3p4/1Pp2r2/1K1R1p1k/8/4P1P1/8 w - c6 0 1",
			[]Move{
				newMove(WhiteKing, NoPiece, B4, A3, NoSquare, NoPiece),
				newMove(WhiteKing, NoPiece, B4, A4, NoSquare, NoPiece),
				newMove(WhiteKing, NoPiece, B4, A5, NoSquare, NoPiece),
				newMove(WhiteKing, NoPiece, B4, B3, NoSquare, NoPiece),
				newMove(WhiteKing, NoPiece, B4, C3, NoSquare, NoPiece),
				newMove(WhiteKing, NoPiece, B4, C4, NoSquare, NoPiece),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			pos := unsafeFEN(tt.args)
			assert.ElementsMatch(t, tt.want, checkFlightMoves(pos))
		})
	}
}

func BenchmarkCheckFlightMoves(b *testing.B) {
	fens := []string{
		"2r2q1k/5pp1/4p1N1/8/1bp5/5P1R/6P1/2R4K b - - 0 1",
		"rnbk1b1r/pp3ppp/2p5/4q1B1/4n3/8/PPP2PPP/2KR1BNR b - - 0 1",
		"8/1Kr5/8/8/8/8/6k1/8 w - - 0 1",
		"8/8/3p4/1Pp2r2/1K1R1p1k/8/4P1P1/8 w - c6 0 1",
	}

	for _, fen := range fens {
		pos := unsafeFEN(fen)
		b.Run(fen, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				checkFlightMoves(pos)
			}
		})
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
			assert.ElementsMatch(t, tt.want, moves)
		})
	}
}

func BenchmarkCastlingMoves(b *testing.B) {
	for _, bb := range testPositions {
		pos := unsafeFEN(bb.preFEN)
		b.Run(bb.preFEN, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				castlingMoves(pos)
			}
		})
	}
}

func TestStandardMoves(t *testing.T) {
	tests := []struct {
		args string
		want []string
	}{
		{
			"1k2q3/8/8/8/8/8/4R3/4K3 w - - 0 1",
			[]string{
				"e1d1", "e1d2", "e1f1", "e1f2", "e2e3",
				"e2e4", "e2e5", "e2e6", "e2e7", "e2e8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			var moves []string
			for _, m := range standardMoves(unsafeFEN(tt.args)) {
				moves = append(moves, m.String())
			}
			assert.ElementsMatch(t, tt.want, moves)
		})
	}
}

func BenchmarkStandardMoves(b *testing.B) {
	for _, bb := range testPositions {
		pos := unsafeFEN(bb.preFEN)
		b.Run(bb.preFEN, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				standardMoves(pos)
			}
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

func BenchmarkMoveBitboard(b *testing.B) {
	benchmarks := []struct {
		sq Square
		pt PieceType
	}{
		{A1, King},
		{B1, Queen},
		{C1, Rook},
		{D1, Bishop},
		{E1, Knight},
		{F1, Pawn},
	}
	pos := unsafeFEN("k7/8/8/8/8/8/5P2/KQRBN3 w - - 0 1")

	for _, bb := range benchmarks {
		b.Run(bb.pt.String(), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				moveBitboard(bb.sq, pos, bb.pt)
			}
		})
	}
}

func TestPawnPushesBitboards(t *testing.T) {
	fenWhite := "k7/p7/1p6/2N5/2n2pP1/1P6/P7/K7 w - - 0 1"
	fenBlack := "k7/p7/1p6/2N5/2n2pP1/1P6/P7/K7 b - g3 0 1"
	posWhite, posBlack := unsafeFEN(fenWhite), unsafeFEN(fenBlack)
	tests := []struct {
		sq  Square
		set []Square
		pos *Position
	}{
		{A2, []Square{A3, A4}, posWhite},
		{B3, []Square{B4}, posWhite},
		{A7, []Square{A5, A6}, posBlack},
		{B6, []Square{B5}, posBlack},
		{F4, []Square{F3}, posBlack},
	}

	for _, tt := range tests {
		t.Run(tt.sq.String(), func(t *testing.T) {
			assert.ElementsMatch(t, tt.set, pawnPushesBitboard(tt.sq, tt.pos).mapping())
		})
	}
}

func TestPawnCapturesBitboards(t *testing.T) {
	fenWhite := "k7/p7/1p6/2N5/2n2pP1/1P6/P7/K7 w - - 0 1"
	fenBlack := "k7/p7/1p6/2N5/2n2pP1/1P6/P7/K7 b - g3 0 1"
	posWhite, posBlack := unsafeFEN(fenWhite), unsafeFEN(fenBlack)
	tests := []struct {
		sq  Square
		set []Square
		pos *Position
	}{
		{A2, nil, posWhite},
		{B3, []Square{C4}, posWhite},
		{A7, nil, posBlack},
		{B6, []Square{C5}, posBlack},
		{F4, []Square{G3}, posBlack},
	}

	for _, tt := range tests {
		t.Run(tt.sq.String(), func(t *testing.T) {
			assert.ElementsMatch(t, tt.set, pawnCapturesBitboard(tt.sq, tt.pos).mapping())
		})
	}
}

func TestCheckBitboard(t *testing.T) {
	pos := unsafeFEN("k6q/8/8/8/8/8/8/K7 w - - 0 1")
	cb := checkBitboard(A1, White, pos.bbOccupied,
		pos.bbBlack&pos.bbKing, pos.bbBlack&pos.bbQueen, pos.bbBlack&pos.bbRook,
		pos.bbBlack&pos.bbBishop, pos.bbBlack&pos.bbKnight, pos.bbBlack&pos.bbPawn)
	assert.Equal(t, newBitboard([]Square{H8}), cb)
}

func BenchmarkCheckBitboard(b *testing.B) {
	pos := unsafeFEN("k6q/8/8/8/8/8/8/K7 w - - 0 1")
	for n := 0; n < b.N; n++ {
		checkBitboard(A1, White, pos.bbOccupied,
			pos.bbBlack&pos.bbKing, pos.bbBlack&pos.bbQueen, pos.bbBlack&pos.bbRook,
			pos.bbBlack&pos.bbBishop, pos.bbBlack&pos.bbKnight, pos.bbBlack&pos.bbPawn)
	}
}

func TestDiagonalBitboard(t *testing.T) {
	occupied := []Square{F6, B6}
	want := []Square{
		B6, F6, C5, E5, C3,
		E3, B2, F2, A1, G1,
	}
	assert.ElementsMatch(t, want, bishopAttacksBitboard(D4, newBitboard(occupied)).mapping())
}

func TestHVBitboard(t *testing.T) {
	occupied := []Square{D3, F5}
	want := []Square{
		D8, D7, D6, A5, B5,
		C5, E5, F5, D4, D3,
	}
	assert.ElementsMatch(t, want, rookAttacksBitboard(D5, newBitboard(occupied)).mapping())
}
