// Package transposition implements different transposition table strategies
// to memoize search results.
package transposition

import (
	"fmt"

	"github.com/notnil/chess"
)

// Entry holds an entry in the transposition hash table.
type Entry struct {
	Score int
	Depth int
	Flag  Flag
}

// Flag represents the score bounds for this entry.
type Flag int8

const (
	NoBounds   Flag = iota // NoBounds represents a score with undefined bounds.
	LowerBound             // LowerBound represents a lower bound score.
	UpperBound             // UpperBound represents an upper bound score.
	Exact                  // Exact represents an exact score.
)

// Interface is the interface implemented by objects that can memoize search results.
type Interface interface {
	fmt.Stringer
	Init(size int) error                   // Init initializes the transposition hash table.
	Set(key *chess.Position, value Entry)  // Set adds an entry to the cache for the given position. If an entry already exists for the position, it is replaced. The addition is not guaranteed.
	Get(key *chess.Position) (Entry, bool) // Get returns the entry (if any) and a boolean representing whether the value was found or not.
	Close()                                // Close initiates a graceful shutdown of the transposition table.
}
