package search

import (
	"context"

	"github.com/notnil/chess"
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
			Position:    input.Position,
			SearchMoves: input.SearchMoves,
			Depth:       depth,
			Evaluation:  input.Evaluation,
		}, -mateScore, mateScore)
		if err != nil {
			return
		}
		o.Mate = matePlies(o.Score)
		output <- o
	}
}

// alphaBeta is the recursive function that implements the Negamax algorithm
// with alpha beta pruning.
func alphaBeta(ctx context.Context, input Input, alpha, beta int) (*Output, error) {
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	default:
	}

	output := terminalNode(input.Position)
	if output != nil {
		return output, nil
	}

	if input.Depth == 0 {
		return &Output{
			Nodes: 1,
			Score: input.Evaluation.Evaluate(input.Position),
		}, nil
	}

	result := &Output{
		Depth: input.Depth,
		Nodes: 0,
		Score: -mateScore,
	}

	for _, move := range searchMoves(input) {
		current, err := alphaBeta(ctx, Input{
			Position:   input.Position.Update(move),
			Depth:      input.Depth - 1,
			Evaluation: input.Evaluation,
		}, -beta, -alpha)
		if err != nil {
			return nil, err
		}

		current.Score = -current.Score
		if current.Score > result.Score {
			result.Score = current.Score
			result.PV = append([]*chess.Move{move}, current.PV...)
		}
		result.Nodes += current.Nodes

		if current.Score > alpha {
			alpha = current.Score
		}

		if alpha >= beta {
			break
		}
	}

	result.Score = updateScore(result.Score)
	return result, nil
}
