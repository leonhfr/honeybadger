package quiescence

import (
	"context"

	"github.com/leonhfr/honeybadger/evaluation"
	"github.com/leonhfr/honeybadger/transposition"
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

	alphaOriginal := input.Alpha

	entry, cached := input.Transposition.Get(input.Position)
	if cached && entry.Depth >= input.Depth {
		switch {
		case entry.Flag == transposition.Exact:
			return &Output{
				Nodes: 1,
				Score: entry.Score,
			}, nil
		case entry.Flag == transposition.LowerBound && entry.Score > input.Alpha:
			input.Alpha = entry.Score
		case entry.Flag == transposition.UpperBound && entry.Score < input.Beta:
			input.Beta = entry.Score
		}

		if input.Alpha >= input.Beta {
			return &Output{
				Nodes: 1,
				Score: entry.Score,
			}, nil
		}
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

	moves := loudMoves(input.Position)
	input.Oracle.Order(moves)

	for _, move := range moves {
		current, err := alphaBeta(ctx, Input{
			Position:      input.Position.Update(move),
			Depth:         input.Depth - 1,
			Alpha:         -input.Beta,
			Beta:          -input.Alpha,
			Evaluation:    input.Evaluation,
			Oracle:        input.Oracle,
			Transposition: input.Transposition,
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

	flag := transposition.Exact
	switch {
	case result.Score <= alphaOriginal:
		flag = transposition.UpperBound
	case result.Score >= input.Beta:
		flag = transposition.LowerBound
	}
	input.Transposition.Set(input.Position, transposition.Entry{
		Score: result.Score,
		Depth: input.Depth,
		Flag:  flag,
	})

	return result, nil
}
