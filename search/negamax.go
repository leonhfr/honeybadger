package search

import (
	"context"

	"github.com/notnil/chess"

	"github.com/leonhfr/honeybadger/evaluation"
)

// Negamax is a variant form of minimax that relies on the
// zero-sum property of a two-player game.
type Negamax struct{}

// String implements the Interface interface.
func (Negamax) String() string {
	return "Negamax"
}

// Search implements the Interface interface.
func (Negamax) Search(ctx context.Context, input Input, output chan<- *Output) {
	for depth := 1; depth <= input.Depth; depth++ {
		o, err := negamax(ctx, Input{
			Position:    input.Position,
			SearchMoves: input.SearchMoves,
			Depth:       depth,
			Evaluation:  input.Evaluation,
		})
		if err != nil {
			return
		}
		o.Mate = mateIn(o.Score)
		output <- o
	}
}

// negamax is the recursive function that implements the Negamax algorithm.
func negamax(ctx context.Context, input Input) (*Output, error) {
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
		Score: -evaluation.Mate,
	}

	for _, move := range searchMoves(input) {
		current, err := negamax(ctx, Input{
			Position:   input.Position.Update(move),
			Depth:      input.Depth - 1,
			Evaluation: input.Evaluation,
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
	}

	result.Score = updateScore(result.Score)
	return result, nil
}

// sign returns the sign +/-1 of the passed integer.
func sign(n int) int {
	if n < 0 {
		return -1
	}
	return 1
}

// updateScore update the score to account for the distance to mate.
func updateScore(score int) int {
	sign := sign(score)
	delta := evaluation.Mate - sign*score
	if delta <= maxDepth {
		return score - sign
	}
	return score
}

// mateIn returns the number of moves before mate.
func mateIn(score int) int {
	sign := sign(score)
	delta := evaluation.Mate - sign*score
	if delta <= maxDepth {
		return sign * (delta/2 + delta%2)
	}
	return 0
}

// searchMoves returns the list of moves to search.
func searchMoves(input Input) []*chess.Move {
	if input.SearchMoves != nil {
		return input.SearchMoves
	}
	return input.Position.ValidMoves()
}

// terminalNode checks if a position is terminal and returns the appropriate score.
func terminalNode(position *chess.Position) *Output {
	switch position.Status() {
	case chess.Checkmate:
		return &Output{
			Nodes: 1,
			Score: -evaluation.Mate, // current player is in checkmate
		}
	case chess.Stalemate,
		chess.ThreefoldRepetition,
		chess.FivefoldRepetition,
		chess.FiftyMoveRule,
		chess.SeventyFiveMoveRule,
		chess.InsufficientMaterial:
		return &Output{
			Nodes: 1,
			Score: evaluation.Draw,
		}
	default:
		return nil
	}
}
