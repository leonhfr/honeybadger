package search

import (
	"context"
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
		select {
		case <-ctx.Done():
			return
		default:
		}

		output <- negamax(Input{
			Position:   input.Position,
			Depth:      depth,
			Evaluation: input.Evaluation,
		})
	}
}

func negamax(input Input) Output {
	if input.Depth == 0 || input.Position.Status() > 0 {
		return Output{
			Nodes: 1,
			Score: input.Evaluation.Evaluate(input.Position),
		}
	}

	result := Output{
		Depth: input.Depth,
		Nodes: 0,
		Score: math.MinInt,
	}

	for _, move := range input.Position.ValidMoves() {
		current := negamax(Input{
			Position:   input.Position.Update(move),
			Depth:      input.Depth - 1,
			Evaluation: input.Evaluation,
		})

		current.Score = -current.Score
		if current.Score > result.Score {
			result.Score = current.Score
			result.PV = append([]*chess.Move{move}, current.PV...)
		}
		result.Nodes += current.Nodes
	}

	return result
}
