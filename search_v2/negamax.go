package search

import (
	"context"

	"github.com/leonhfr/honeybadger/chess"
)

func negamax(ctx context.Context, pos *chess.Position, depth int) (*output, error) {
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	default:
	}

	moves := pos.PseudoMoves()
	score, terminal := isTerminal(pos, len(moves), depth)
	if terminal {
		return &output{
			nodes: 1,
			score: score,
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
