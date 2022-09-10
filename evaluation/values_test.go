package evaluation

import (
	"testing"

	"github.com/notnil/chess"
	"github.com/stretchr/testify/assert"
)

func TestPieceValuesEvaluate(t *testing.T) {
	strategy := Values{}
	tests := []struct {
		name string
		args string
		want int
	}{
		{name: "starting position", args: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fen, err := chess.FEN(tt.args)
			assert.NoErrorf(t, err, "could not parse FEN %s", tt.args)

			game := chess.NewGame(fen)
			assert.Equal(t, tt.want, strategy.Evaluate(game.Position()))
		})
	}
}
