package search

import (
	"context"
	"math/rand"

	"github.com/notnil/chess"
)

// Capture tries to capture enemy pieces and otherwise plays random moves.
type Capture struct{}

// String implements the Interface interface.
func (Capture) String() string {
	return "Capture"
}

// Search implements the Interface interface.
func (Capture) Search(ctx context.Context, input Input, output chan<- Output) {
	moves := input.Position.ValidMoves()

	var captures []*chess.Move
	for _, m := range moves {
		if m.HasTag(chess.Capture) {
			captures = append(captures, m)
		}
	}

	if len(captures) > 0 {
		pv := []*chess.Move{captures[rand.Intn(len(captures))]} //nolint
		output <- Output{
			PV: pv,
		}
		return
	}

	pv := []*chess.Move{moves[rand.Intn(len(moves))]} //nolint
	output <- Output{
		PV: pv,
	}
}
