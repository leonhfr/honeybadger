package search

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	for _, tt := range testCheckmatePositions {
		t.Run(tt.name, func(t *testing.T) {
			pos := unsafeFEN(tt.fen)
			output, err := search(context.Background(), pos, -mate, mate, tt.depth)

			assert.Equal(t, tt.output.score, output.score)
			assert.Equal(t, tt.moves, movesString(output.pv))
			assert.Nil(t, err)
		})
	}
}

func BenchmarkSearch(b *testing.B) {
	for _, bb := range testCheckmatePositions {
		b.Run(bb.name, func(b *testing.B) {
			pos := unsafeFEN(bb.fen)
			for n := 0; n < b.N; n++ {
				_, _ = search(context.Background(), pos, -mate, mate, bb.depth)
			}
		})
	}
}
