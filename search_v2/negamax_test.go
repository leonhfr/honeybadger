package search

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/honeybadger/chess"
)

var testCheckmatePositions = []struct {
	name   string
	fen    string
	depth  int
	output output
	moves  []string
}{
	{
		name:   "checkmate",
		fen:    "8/8/8/5K1k/8/8/8/7R b - - 0 1",
		depth:  1,
		output: output{0, 1, -mate, nil},
		moves:  nil,
	},
	{
		name:   "mate in 1",
		fen:    "8/8/8/5K1k/8/8/8/5R2 w - - 0 1",
		depth:  1,
		output: output{1, 15, mate - 1, nil},
		moves:  []string{"f1h1"},
	},
	{
		name:   "mate in 1",
		fen:    "r1b1kb1r/pppp1ppp/2n1pq2/8/3Pn2N/2P3P1/PP1NPP1P/R1BQKB1R b KQkq - 3 6",
		depth:  1,
		output: output{1, 46, mate - 1, nil},
		moves:  []string{"f6f2"},
	},
	{
		name:   "mate in 2",
		fen:    "5rk1/pb2npp1/1pq4p/5p2/5B2/1B6/P2RQ1PP/2r1R2K b - - 0 1",
		depth:  3,
		output: output{3, 90094, mate - 3, nil},
		moves:  []string{"c1e1", "e2g2", "c6g2"},
	},
}

func negamax(ctx context.Context, pos *chess.Position, depth int) (*output, error) {
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	default:
	}

	moves := pos.PseudoMoves()
	score, terminal := isTerminal(pos, len(moves))
	if terminal {
		return &output{
			nodes: 1,
			score: score,
		}, nil
	}

	if depth == 0 {
		return &output{
			nodes: 1,
			score: evaluate(pos),
		}, nil
	}

	result := &output{
		depth: depth,
		nodes: 0,
		score: -mate,
	}

	for _, move := range moves {
		metadata, ok := pos.MakeMove(move)
		if !ok {
			continue
		}

		current, err := negamax(ctx, pos, depth-1)
		if err != nil {
			return nil, err
		}

		current.score = -current.score
		if current.score > result.score {
			result.score = current.score
			result.pv = append(current.pv, move)
		}
		result.nodes += current.nodes

		pos.UnmakeMove(move, metadata)
	}

	result.score = incMateDistance(result.score)
	return result, nil
}

func TestNegamax(t *testing.T) {
	for _, tt := range testCheckmatePositions {
		t.Run(tt.name, func(t *testing.T) {
			pos := unsafeFEN(tt.fen)
			output, err := negamax(context.Background(), pos, tt.depth)

			assert.Equal(t, tt.output.nodes, output.nodes)
			assert.Equal(t, tt.output.score, output.score)
			assert.Equal(t, tt.moves, movesString(output.pv))
			assert.Nil(t, err)
		})
	}
}

func BenchmarkNegamax(b *testing.B) {
	for _, bb := range testCheckmatePositions {
		b.Run(bb.name, func(b *testing.B) {
			pos := unsafeFEN(bb.fen)
			for n := 0; n < b.N; n++ {
				_, _ = negamax(context.Background(), pos, bb.depth)
			}
		})
	}
}

func movesString(moves []chess.Move) []string {
	var result []string
	for _, move := range moves {
		result = append(result, move.String())
	}
	return result
}
