package search

import (
	"math"

	"github.com/leonhfr/honeybadger/chess"
)

const (
	// maxDepth is the maximum depth at which the package will search.
	maxDepth = 64
	// mate is the score of a checkmate.
	mate = math.MaxInt
	// draw is the score of a draw.
	draw = 0
)

// output holds a search output.
type output struct {
	depth int
	nodes int
	score int
	pv    []chess.Move // reversed
}
