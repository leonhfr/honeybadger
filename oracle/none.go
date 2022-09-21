package oracle

import "github.com/notnil/chess"

// None is the strategy used when we want no move ordering.
type None struct{}

// String implements the Interface interface.
func (None) String() string {
	return "None"
}

// Order implements the Interface interface.
func (None) Order(moves []*chess.Move) {}
