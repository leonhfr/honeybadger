package search

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/honeybadger/chess"
)

func alphaBeta(ctx context.Context, pos *chess.Position, alpha, beta, depth int) (*output, error) {
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

	orderMoves(moves)
	for _, move := range moves {
		metadata, ok := pos.MakeMove(move)
		if !ok {
			continue
		}

		current, err := alphaBeta(ctx, pos, -beta, -alpha, depth-1)
		if err != nil {
			return nil, err
		}

		current.score = -current.score
		if current.score > result.score {
			result.score = current.score
			result.pv = append(current.pv, move)
		}
		result.nodes += current.nodes

		if current.score > alpha {
			alpha = current.score
		}

		pos.UnmakeMove(move, metadata)

		if alpha >= beta {
			break
		}
	}

	result.score = incMateDistance(result.score)
	return result, nil
}

func TestAlphaBeta(t *testing.T) {
	for _, tt := range testCheckmatePositions {
		t.Run(tt.name, func(t *testing.T) {
			pos := unsafeFEN(tt.fen)
			output, err := alphaBeta(context.Background(), pos, -mate, mate, tt.depth)

			assert.Equal(t, tt.output.score, output.score)
			assert.Equal(t, tt.moves, movesString(output.pv))
			assert.Nil(t, err)
		})
	}
}

func BenchmarkAlphaBeta(b *testing.B) {
	for _, bb := range testCheckmatePositions {
		b.Run(bb.name, func(b *testing.B) {
			pos := unsafeFEN(bb.fen)
			for n := 0; n < b.N; n++ {
				_, _ = alphaBeta(context.Background(), pos, -mate, mate, bb.depth)
			}
		})
	}
}
