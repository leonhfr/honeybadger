// Package oracle implements move ordering. Moves that have better
// chances to result in pruning are sorted first.
package oracle

import (
	"fmt"

	"github.com/notnil/chess"
)

// Interface is the interface implemented by objects that can order moves.
type Interface interface {
	fmt.Stringer
	// Order sorts the moves. The moves that have better chances to result in
	// earlier and more frequent alpha/beta cut offs are sorter first.
	Order(moves []*chess.Move)
}
