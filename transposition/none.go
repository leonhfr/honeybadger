package transposition

import "github.com/notnil/chess"

// None is the strategy used when we want no transposition hash tables.
type None struct{}

// String implements the Interface interface.
func (None) String() string {
	return "None"
}

// Init implements the Interface interface.
func (None) Init() error {
	return nil
}

// Set implements the Interface interface.
func (None) Set(key *chess.Position, value Entry) {}

// Get implements the Interface interface.
func (None) Get(key *chess.Position) (Entry, bool) {
	return Entry{}, false
}

// Close implements the Interface interface.
func (None) Close() {}
