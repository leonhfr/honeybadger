package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderMoves(t *testing.T) {
	tests := []struct {
		name string
		args string
		want []string
	}{
		{
			"promotions",
			"7k/P7/8/8/8/8/8/K7 w - - 0 1",
			[]string{"a7a8q", "a7a8n", "a1b1", "a1a2", "a1b2", "a7a8r", "a7a8b"},
		},
		{
			"tags",
			"rnbq1knr/p1pp2pp/8/Pp6/8/8/8/R3K2R w KQ b6 0 1",
			[]string{
				"e1g1", "e1c1", "h1h7", "a5b6", "a1a2", "h1h5", "a1a3", "a1a4",
				"h1f1", "h1g1", "h1h2", "h1h3", "h1h4", "a5a6", "h1h6", "a1d1",
				"e1d1", "e1f1", "e1d2", "e1e2", "e1f2", "a1c1", "a1b1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moves := unsafeFEN(tt.args).PseudoMoves()
			orderMoves(moves)

			var got []string
			for _, move := range moves {
				got = append(got, move.String())
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
