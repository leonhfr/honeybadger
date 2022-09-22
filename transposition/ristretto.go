package transposition

import (
	"unsafe"

	"github.com/dgraph-io/ristretto"
	"github.com/notnil/chess"
)

// Ristretto implements transposition hash tables using the Ristretto library.
// It is safe to use concurrently.
type Ristretto struct {
	cache *ristretto.Cache
}

// String implements the Interface interface.
func (Ristretto) String() string {
	return "Ristretto"
}

// Init implements the Interface interface.
func (r *Ristretto) Init(size int) error {
	bytes := uint64(unsafe.Sizeof(Entry{}))
	maxCost := int64(1024 * 1024 * uint64(size) / bytes)

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 10 * maxCost,
		MaxCost:     maxCost,
		BufferItems: 64,
	})
	if err != nil {
		return err
	}

	r.cache = cache
	return nil
}

// Set implements the Interface interface.
func (r *Ristretto) Set(key *chess.Position, entry Entry) {
	r.cache.Set(key.Hash(), entry, 1)
}

// Get implements the Interface interface.
func (r *Ristretto) Get(key *chess.Position) (Entry, bool) {
	entry, found := r.cache.Get(key.Hash())
	if !found {
		return Entry{}, false
	}
	return entry.(Entry), true
}

// Close implements the Interface interface.
func (r *Ristretto) Close() {
	r.cache.Close()
}
