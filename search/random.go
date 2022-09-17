package search

import (
	"context"
	"math/rand"

	"github.com/notnil/chess"
)

// Random plays random moves.
type Random struct{}

// String implements the Interface interface.
func (Random) String() string {
	return "Random"
}

// Search implements the Interface interface.
func (Random) Search(ctx context.Context, input Input, output chan<- *Output) {
	moves := input.Position.ValidMoves()
	pv := []*chess.Move{moves[rand.Intn(len(moves))]} //nolint
	output <- &Output{
		PV: pv,
	}
}
