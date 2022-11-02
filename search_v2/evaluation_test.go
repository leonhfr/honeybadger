package search

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/honeybadger/chess"
)

var testPositions = []struct {
	fen   string
	score int
}{
	{
		fen:   "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		score: 0,
	},
	{
		fen:   "2r3k1/1q1nbppp/r3p3/3pP3/pPpP4/P1Q2N2/2RN1PPP/2R4K b - b3 0 23",
		score: 10,
	},
	{
		fen:   "r2qk2r/pp1n1ppp/2pbpn2/3p4/2PP4/1PNQPN2/P4PPP/R1B1K2R w KQkq - 1 9",
		score: -25,
	},
	{
		fen:   "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		score: 0,
	},
	{
		fen:   "r3k2r/ppqn1ppp/2pbpn2/3p4/2PP4/1PNQPN2/P2B1PPP/R3K2R w KQkq - 3 10",
		score: -13,
	},
	{
		fen:   "r1bqkbnr/ppp1pppp/2n5/3p4/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3",
		score: 0,
	},
	{
		fen:   "r1bqkbnr/ppp1p1pp/2n5/3pPp2/8/5N2/PPPP1PPP/RNBQKB1R w KQkq f6 0 4",
		score: 24,
	},
	{
		fen:   "r1bqkbnr/ppp1p1pp/2n5/3pPp2/3N4/8/PPPP1PPP/RNBQKB1R b KQkq - 1 4",
		score: -20,
	},
	{
		fen:   "r7/1Pp5/2P3p1/8/6pb/4p1kB/4P1p1/6K1 w - - 0 1",
		score: -723,
	},
}

func TestEvaluate(t *testing.T) {
	for _, tt := range testPositions {
		t.Run(tt.fen, func(t *testing.T) {
			pos := unsafeFEN(tt.fen)
			assert.Equal(t, tt.score, evaluate(pos))
		})
	}
}

func BenchmarkEvaluate(b *testing.B) {
	for _, bb := range testPositions {
		b.Run(bb.fen, func(b *testing.B) {
			pos := unsafeFEN(bb.fen)
			for n := 0; n < b.N; n++ {
				evaluate(pos)
			}
		})
	}
}

func unsafeFEN(fen string) *chess.Position {
	p, err := chess.FromFEN(fen)
	if err != nil {
		panic(err)
	}
	return p
}
