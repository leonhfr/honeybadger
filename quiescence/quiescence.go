// Package quiescence implements different quiescence search strategies.
package quiescence

import (
	"context"
	"fmt"

	"github.com/notnil/chess"

	"github.com/leonhfr/honeybadger/evaluation"
)

// Input holds a quiescence search input.
type Input struct {
	Position   *chess.Position      // Current board position.
	Depth      int                  // Search <x> plies only.
	Alpha      int                  // Best score that the maximizer can guarantee.
	Beta       int                  // Best score that the minimizer can guarantee.
	Evaluation evaluation.Interface // Evaluation strategy to use.
}

// Output holds a quiescence search output.
type Output struct {
	Nodes int // Number of nodes searched.
	Score int // Score from the engine's point of view in centipawns.
}

// Interface is the interface implemented by objects that can
// run a quiescence search on a chess board.
//
// The Search method includes alpha-beta integers arguments for the cases
// the algorithm implements alpha-beta pruning.
type Interface interface {
	fmt.Stringer
	Search(ctx context.Context, input Input) (*Output, error)
}

const (
	// MaxDepth is the maximum depth at which the quiescence package will search.
	MaxDepth = 4
)

// IsQuiet determines whether a position is quiet.
//
// A position is quiet when it has no loud moves.
func IsQuiet(position *chess.Position) bool {
	return len(loudMoves(position)) == 0
}

// loudMoves returns the list of loud moves from a position.
//
// A loud move is a move that captures another pieces.
func loudMoves(position *chess.Position) []*chess.Move {
	var moves []*chess.Move
	for _, m := range position.ValidMoves() {
		if m.HasTag(chess.Capture) || m.HasTag(chess.EnPassant) {
			moves = append(moves, m)
		}
	}
	return moves
}
