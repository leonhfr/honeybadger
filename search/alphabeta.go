package search

import (
	"context"

	"github.com/notnil/chess"

	"github.com/leonhfr/honeybadger/evaluation"
	"github.com/leonhfr/honeybadger/quiescence"
	"github.com/leonhfr/honeybadger/transposition"
)

// AlphaBeta pruning is an optimization for Negamax. It returns the same
// result as Negamax but decreases the number of nodes to evaluate.
type AlphaBeta struct{}

// String implements the Interface interface.
func (AlphaBeta) String() string {
	return "AlphaBeta"
}

// Search implements the Interface interface.
func (AlphaBeta) Search(ctx context.Context, input Input, output chan<- *Output) {
	for depth := 1; depth <= input.Depth; depth++ {
		o, err := alphaBeta(ctx, Input{
			Position:      input.Position,
			SearchMoves:   input.SearchMoves,
			Depth:         depth,
			alpha:         -evaluation.Mate,
			beta:          evaluation.Mate,
			Evaluation:    input.Evaluation,
			Oracle:        input.Oracle,
			Quiescence:    input.Quiescence,
			Transposition: input.Transposition,
		})
		if err != nil {
			return
		}
		o.Mate = mateIn(o.Score)
		output <- o
	}
}

// alphaBeta is the recursive function that implements the Negamax algorithm
// with alpha beta pruning.
func alphaBeta(ctx context.Context, input Input) (*Output, error) {
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	default:
	}

	alphaOriginal := input.alpha

	entry, cached := input.Transposition.Get(input.Position)
	if cached && entry.Depth >= input.Depth {
		switch {
		case entry.Flag == transposition.Exact:
			return &Output{
				Nodes: 1,
				Score: entry.Score,
			}, nil
		case entry.Flag == transposition.LowerBound && entry.Score > input.alpha:
			input.alpha = entry.Score
		case entry.Flag == transposition.UpperBound && entry.Score < input.beta:
			input.beta = entry.Score
		}

		if input.alpha >= input.beta {
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

	if input.Depth == 0 {
		if quiescence.IsQuiet(input.Position) {
			return &Output{
				Nodes: 1,
				Score: input.Evaluation.Evaluate(input.Position),
			}, nil
		}

		output, err := input.Quiescence.Search(ctx, quiescence.Input{
			Position:      input.Position,
			Depth:         quiescence.MaxDepth,
			Alpha:         -input.beta,
			Beta:          -input.alpha,
			Evaluation:    input.Evaluation,
			Oracle:        input.Oracle,
			Transposition: input.Transposition,
		})
		if err != nil {
			return nil, err
		}

		return &Output{
			Nodes: output.Nodes,
			Score: output.Score,
		}, nil
	}

	result := &Output{
		Depth: input.Depth,
		Nodes: 0,
		Score: -evaluation.Mate,
	}

	moves := searchMoves(input)
	input.Oracle.Order(moves)

	for _, move := range moves {
		current, err := alphaBeta(ctx, Input{
			Position:      input.Position.Update(move),
			Depth:         input.Depth - 1,
			alpha:         -input.beta,
			beta:          -input.alpha,
			Evaluation:    input.Evaluation,
			Oracle:        input.Oracle,
			Quiescence:    input.Quiescence,
			Transposition: input.Transposition,
		})
		if err != nil {
			return nil, err
		}

		current.Score = -current.Score
		if current.Score > result.Score {
			result.Score = current.Score
			result.PV = append([]*chess.Move{move}, current.PV...)
		}
		result.Nodes += current.Nodes

		if current.Score > input.alpha {
			input.alpha = current.Score
		}

		if input.alpha >= input.beta {
			break
		}
	}

	result.Score = evaluation.IncMateDistance(result.Score, maxDepth)

	flag := transposition.Exact
	switch {
	case result.Score <= alphaOriginal:
		flag = transposition.UpperBound
	case result.Score >= input.beta:
		flag = transposition.LowerBound
	}
	input.Transposition.Set(input.Position, transposition.Entry{
		Score: result.Score,
		Depth: input.Depth,
		Flag:  flag,
	})

	return result, nil
}
