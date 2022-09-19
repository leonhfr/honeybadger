package quiescence

import (
	"context"

	"github.com/leonhfr/honeybadger/evaluation"
)

// AlphaBeta performs a quiescence search using the negamax search algorithm
// and alpha-beta pruning.
type AlphaBeta struct{}

// String implements the Interface interface.
func (AlphaBeta) String() string {
	return "AlphaBeta"
}

// Search implements the Interface interface.
func (AlphaBeta) Search(ctx context.Context, input Input) (*Output, error) {
	return alphaBeta(ctx, input)
}

// alphaBeta is the recursive function that implements quiescence search using the
// negamax algorithm with alpha-beta pruning.
func alphaBeta(ctx context.Context, input Input) (*Output, error) {
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	default:
	}

	score, terminal := evaluation.Terminal(input.Position)
	if terminal {
		return &Output{
			Nodes: 1,
			Score: score,
		}, nil
	}

	if input.Depth == 0 || IsQuiet(input.Position) {
		return &Output{
			Nodes: 1,
			Score: input.Evaluation.Evaluate(input.Position),
		}, nil
	}

	result := &Output{
		Nodes: 0,
		Score: -evaluation.Mate,
	}

	for _, move := range loudMoves(input.Position) {
		current, err := alphaBeta(ctx, Input{
			Position:   input.Position.Update(move),
			Depth:      input.Depth - 1,
			Evaluation: input.Evaluation,
			Alpha:      -input.Beta,
			Beta:       -input.Alpha,
		})
		if err != nil {
			return nil, err
		}

		current.Score = -current.Score
		if current.Score > result.Score {
			result.Score = current.Score
		}
		result.Nodes += current.Nodes

		if current.Score > input.Alpha {
			input.Alpha = current.Score
		}

		if input.Alpha >= input.Beta {
			break
		}
	}

	result.Score = evaluation.IncMateDistance(result.Score, MaxDepth)
	return result, nil
}
