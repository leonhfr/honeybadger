// Package search implements different search strategies.
package search

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/notnil/chess"
)

// Input holds a search input.
type Input struct {
	Position *chess.Position // Current board position.
	Strategy Interface       // Search strategy to use.
}

// Output holds a search output.
type Output struct {
	Depth int           // Search depth in plies.
	Nodes int           // Number of nodes searched.
	Score int           // Score from the engine's point of view in centipawns.
	PV    []*chess.Move // Principal variation, best line found.
}

// Interface is the interface implemented by objects that can
// run a search on a chess board.
type Interface interface {
	fmt.Stringer
	Search(ctx context.Context, input Input, output chan<- Output)
}

// Run starts a search.
func Run(ctx context.Context, input Input) <-chan Output {
	output := make(chan Output)

	go func() {
		defer close(output)
		input.Strategy.Search(ctx, input, output)
	}()

	return output
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
