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
	Evaluate(p *chess.Position) int // Evaluate returns the score of a given chess position.
}

const (
	// Mate is the score of a checkmate.
	Mate = math.MaxInt
	// Draw is the score of a draw.
	Draw = 0
)

// Terminal checks if a position is terminal and returns a tuple (int, bool)
// returning the position score and whether it is terminal.
func Terminal(position *chess.Position) (int, bool) {
	switch position.Status() {
	case chess.Checkmate:
		return -Mate, true
	case chess.Stalemate,
		chess.ThreefoldRepetition,
		chess.FivefoldRepetition,
		chess.FiftyMoveRule,
		chess.SeventyFiveMoveRule,
		chess.InsufficientMaterial:
		return Draw, true
	default:
		return 0, false
	}
}

// IncMateDistance increase the distance to the mate by a count of one.
//
// In case of a positive score, it is decreased by 1.
// In case of a negative score, it is increased by 1.
func IncMateDistance(score, maxDepth int) int {
	sign := 1
	if score < 0 {
		sign = -1
	}
	delta := Mate - sign*score
	if delta <= maxDepth {
		return score - sign
	}
	return score
}
