package search

import (
	"context"

	"github.com/leonhfr/honeybadger/chess"
)

const (
	// maxDepth is the maximum depth at which the package will search.
	maxDepth = 64
)

func negamax(ctx context.Context, pos *chess.Position, depth int) (*output, error) {
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	default:
	}

	moves := pos.PseudoMoves()

	if len(moves) == 0 && pos.InCheck() {
		return &output{
			nodes: 1,
			score: -mate,
		}, nil
	} else if len(moves) == 0 {
		return &output{
			nodes: 1,
			score: draw,
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

	result.score = incMateDistance(result.score, maxDepth)
	return result, nil
}

// output holds a search output.
type output struct {
	depth int
	nodes int
	score int
	pv    []chess.Move // reversed
}
