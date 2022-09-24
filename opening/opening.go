// Package opening provides access to opening books.
package opening

import (
	"bytes"
	"fmt"
	"math/rand"

	"github.com/notnil/chess"

	"github.com/leonhfr/honeybadger/opening/polyglot"
)

// Interface is the interface implemented by objects that can return opening moves.
type Interface interface {
	fmt.Stringer
	// Init initializes the opening book.
	Init(book []byte) error
	// Move returns the move to be played from the opening book.
	// If no move is found, nil is returned.
	Move(position *chess.Position) *chess.Move
}

// None is the strategy used when we want no moves to be played from an
// opening book.
type None struct{}

// NewNone returns a new None opening strategy.
func NewNone() *None {
	return &None{}
}

// String implements the Interface interface.
func (*None) String() string {
	return "None"
}

// Init implements the Interface interface.
func (*None) Init(data []byte) error {
	return nil
}

// Move implements the Interface interface.
func (*None) Move(position *chess.Position) *chess.Move {
	return nil
}

// Best returns the best move found from the opening book.
type Best struct {
	book *polyglot.Book
}

// NewBest returns a new Best opening strategy.
func NewBest() *Best {
	return &Best{polyglot.New()}
}

// String implements the Interface interface.
func (*Best) String() string {
	return "Best"
}

// Init implements the Interface interface.
func (b *Best) Init(data []byte) error {
	r := bytes.NewReader(data)
	return b.book.Init(r)
}

// Move implements the Interface interface.
func (b *Best) Move(position *chess.Position) *chess.Move {
	moves := b.book.Lookup(position)
	if len(moves) == 0 {
		return nil
	}
	return moves[0].Move
}

// UniformRandom returns a random move with uniform probabilities.
type UniformRandom struct {
	book *polyglot.Book
}

// NewUniformRandom returns a new UniformRandom opening strategy.
func NewUniformRandom() *UniformRandom {
	return &UniformRandom{polyglot.New()}
}

// String implements the Interface interface.
func (*UniformRandom) String() string {
	return "UniformRandom"
}

// Init implements the Interface interface.
func (ur *UniformRandom) Init(data []byte) error {
	r := bytes.NewReader(data)
	return ur.book.Init(r)
}

// Move implements the Interface interface.
func (ur *UniformRandom) Move(position *chess.Position) *chess.Move {
	moves := ur.book.Lookup(position)
	if len(moves) == 0 {
		return nil
	}
	index := rand.Intn(len(moves)) //nolint
	return moves[index].Move
}

// WeightedRandom returns a random move with weighted probabilities from the
// opening book.
type WeightedRandom struct {
	book *polyglot.Book
}

// NewWeightedRandom returns a new WeightedRandom opening strategy.
func NewWeightedRandom() *WeightedRandom {
	return &WeightedRandom{polyglot.New()}
}

// String implements the Interface interface.
func (*WeightedRandom) String() string {
	return "WeightedRandom"
}

// Init implements the Interface interface.
func (wr *WeightedRandom) Init(data []byte) error {
	r := bytes.NewReader(data)
	return wr.book.Init(r)
}

// Move implements the Interface interface.
func (wr *WeightedRandom) Move(position *chess.Position) *chess.Move {
	moves := wr.book.Lookup(position)
	if len(moves) == 0 {
		return nil
	}
	var sum int
	for _, move := range moves {
		sum += move.Weight
	}
	index := rand.Intn(sum) //nolint
	for _, move := range moves {
		if index < move.Weight {
			return move.Move
		}
		index -= move.Weight
	}
	return nil
}
