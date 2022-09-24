package polyglot

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	data "github.com/leonhfr/honeybadger/test/data/polyglot"
)

func TestNew(t *testing.T) {
	book := New()
	assert.Equal(t, map[uint64][]entry{}, book.positions)
}

func TestInit(t *testing.T) {
	book := New()
	reader := bytes.NewReader(data.LaskerTrap)
	err := book.Init(reader)

	if ok := assert.Nil(t, err); ok {
		assert.Len(t, book.positions, 7)
	}
}

func TestLookup(t *testing.T) {
	book := New()
	reader := bytes.NewReader(data.LaskerTrap)
	err := book.Init(reader)

	if ok := assert.Nil(t, err); ok {
		pos := position("rnbqk1nr/ppp2ppp/8/4P3/1bPp4/4P3/PP1B1PPP/RN1QKBNR b KQkq - 0 1")
		moves := book.Lookup(pos)
		if ok := assert.Len(t, moves, 1); ok {
			assert.Equal(t, "d4e3", moves[0].Move.String())
			assert.Equal(t, 2, moves[0].Weight)
		}
	}
}
