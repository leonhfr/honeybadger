// Package search implements different search strategies.
package search

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/notnil/chess"

	"github.com/leonhfr/honeybadger/evaluation"
	"github.com/leonhfr/honeybadger/quiescence"
)

// Input holds a search input.
type Input struct {
	Position    *chess.Position      // Current board position.
	SearchMoves []*chess.Move        // Restrict search to those moves only.
	Depth       int                  // Search <x> plies only.
	Search      Interface            // Search strategy to use.
	Evaluation  evaluation.Interface // Evaluation strategy to use.
	Quiescence  quiescence.Interface // Quiescence strategy to use.
}

// Output holds a search output.
type Output struct {
	Depth int           // Search depth in plies.
	Nodes int           // Number of nodes searched.
	Score int           // Score from the engine's point of view in centipawns.
	Mate  int           // Number of moves before mate. Positive for the current player to mate, negative for the current player to be mated.
	PV    []*chess.Move // Principal variation, best line found.
}

// Interface is the interface implemented by objects that can
// run a search on a chess board.
type Interface interface {
	fmt.Stringer
	Search(ctx context.Context, input Input, output chan<- *Output) // Search runs a search.
}

// Run starts a search.
func Run(ctx context.Context, input Input) <-chan *Output {
	output := make(chan *Output)

	if input.Depth == 0 {
		input.Depth = maxDepth
	}

	go func() {
		defer close(output)
		input.Search.Search(ctx, input, output)
	}()

	return output
}

const (
	// maxDepth is the maximum depth at which the package will search.
	maxDepth = 64
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
