package search

import (
	"context"
	"errors"
	"math"

	"github.com/notnil/chess"
)

// Negamax is a variant form of minimax that relies on the
// zero-sum property of a two-player game.
type Negamax struct{}

// String implements the Interface interface.
func (Negamax) String() string {
	return "Negamax"
}

// Search implements the Interface interface.
func (Negamax) Search(ctx context.Context, input Input, output chan<- Output) {
	for depth := 1; depth < input.Depth; depth++ {
		o, err := negamax(ctx, Input{
			Position:   input.Position,
			Depth:      depth,
			Evaluation: input.Evaluation,
		})
		if err != nil {
			return
		}
		output <- o
	}
}

func negamax(ctx context.Context, input Input) (Output, error) {
	select {
	case <-ctx.Done():
		return Output{}, errors.New("aborted by context")
	default:
	}

	if input.Depth == 0 || input.Position.Status() > 0 {
		return Output{
			Nodes: 1,
			Score: input.Evaluation.Evaluate(input.Position),
		}, nil
	}

	result := Output{
		Depth: input.Depth,
		Nodes: 0,
		Score: math.MinInt,
	}

	for _, move := range input.Position.ValidMoves() {
		current, err := negamax(ctx, Input{
			Position:   input.Position.Update(move),
			Depth:      input.Depth - 1,
			Evaluation: input.Evaluation,
		})
		if err != nil {
			return result, err
		}

		current.Score = -current.Score
		if current.Score > result.Score {
			result.Score = current.Score
			result.PV = append([]*chess.Move{move}, current.PV...)
		}
		result.Nodes += current.Nodes
	}

	return result, nil
}
