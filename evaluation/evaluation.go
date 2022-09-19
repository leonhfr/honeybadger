// Package evaluation implements different evaluation strategies.
package evaluation

import (
	"fmt"
	"math"

	"github.com/notnil/chess"
)

// Interface is the interface implemented by objects that can
// evaluate a chess board position. It returns a value from the
// point of view of the position's current player.
type Interface interface {
	fmt.Stringer
	Evaluate(p *chess.Position) int
}

const (
	// Mate is the score of a checkmate.
	Mate = math.MaxInt
	// Draw is the score of a draw.
	Draw = 0
)
