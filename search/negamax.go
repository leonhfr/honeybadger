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
func (Negamax) Search(ctx context.Context, input Input, output chan<- *Output) {
	for depth := 1; depth <= input.Depth; depth++ {
		select {
		case <-ctx.Done():
			return
		default:
		}

		output <- negamax(Input{
			Position:    input.Position,
			SearchMoves: input.SearchMoves,
			Depth:       depth,
			Evaluation:  input.Evaluation,
		})
	}
}

// negamax is the recursive function that implements the Negamax algorithm
func negamax(input Input) *Output {
	output := terminalNode(input.Position)
	if output != nil {
		return output
	}

	if input.Depth == 0 {
		return &Output{
			Nodes: 1,
			Score: input.Evaluation.Evaluate(input.Position),
		}
	}

	result := &Output{
		Depth: input.Depth,
		Nodes: 0,
		Score: math.MinInt,
	}

	for _, move := range searchMoves(input) {
		current := negamax(Input{
			Position:   input.Position.Update(move),
			Depth:      input.Depth - 1,
			Evaluation: input.Evaluation,
		})

		current.Score = -current.Score
		current.Mate = -current.Mate
		updateResult(result, current, move)
		result.Nodes += current.Nodes
	}

	return result
}

// updateResult updates the result if the current output is better
func updateResult(result, current *Output, move *chess.Move) {
	// move that don't lead to mate but is better
	if !current.mate && current.Score > result.Score {
		result.Score = current.Score
		result.PV = append([]*chess.Move{move}, current.PV...)
	}

	// move that lead to mate when we don't have any currently
	// or move that lead to mate that do so in fewer moves
	if current.mate && (!result.mate || current.Mate < result.Mate) {
		result.Score = current.Score
		result.PV = append([]*chess.Move{move}, current.PV...)
		inc := 1
		if result.Score < 0 {
			inc = -1
		}
		result.mate = true
		result.Mate = current.Mate + inc
	}
}

// searchMoves returns the list of moves to search
func searchMoves(input Input) []*chess.Move {
	if input.SearchMoves != nil {
		return input.SearchMoves
	}
	return input.Position.ValidMoves()
}

// terminalNode checks if a position is terminal and returns the appropriate score
func terminalNode(position *chess.Position) *Output {
	switch position.Status() {
	case chess.Checkmate:
		return &Output{
			Nodes: 1,
			Score: math.MinInt + 1, // +1 allows negation to positive score
			mate:  true,
		}
	case chess.Stalemate,
		chess.ThreefoldRepetition,
		chess.FivefoldRepetition,
		chess.FiftyMoveRule,
		chess.SeventyFiveMoveRule,
		chess.InsufficientMaterial:
		return &Output{
			Nodes: 1,
			Score: 0, // draw
		}
	default:
		return nil
	}
}
