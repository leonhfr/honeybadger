package quiescence

import (
	"context"
	"testing"

	"github.com/notnil/chess"

	"github.com/leonhfr/honeybadger/evaluation"
	"github.com/leonhfr/honeybadger/oracle"
	"github.com/leonhfr/honeybadger/transposition"
)

func position(fen string) *chess.Position {
	fn, _ := chess.FEN(fen)
	game := chess.NewGame(fn)
	game.Outcome()
	return game.Position()
}

func benchmarkAlphaBeta(fen string, depth int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = alphaBeta(context.Background(), Input{
			Position:      position(fen),
			Depth:         depth,
			Alpha:         -evaluation.Mate,
			Beta:          evaluation.Mate,
			Evaluation:    evaluation.Pesto{},
			Oracle:        oracle.Order{},
			Transposition: transposition.None{},
		})
	}
}

func BenchmarkAlphaBeta1(b *testing.B) {
	benchmarkAlphaBeta("r1b1kb1r/pppp1ppp/2n1pq2/8/3Pn2N/2P3P1/PP1NPP1P/R1BQKB1R b KQkq - 3 6", 1, b)
}

func BenchmarkAlphaBeta3(b *testing.B) {
	benchmarkAlphaBeta("5rk1/pb2npp1/1pq4p/5p2/5B2/1B6/P2RQ1PP/2r1R2K b - - 0 1", 3, b)
}
