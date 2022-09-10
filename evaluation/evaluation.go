// Package evaluation implements different evaluation strategies.
package evaluation

import (
	"fmt"

	"github.com/notnil/chess"
)

// Interface is the interface implemented by objects that can
// evaluate a chess board position. It returns a value from the
// point of view of the position's current player.
type Interface interface {
	fmt.Stringer
	Evaluate(p *chess.Position) int
}
