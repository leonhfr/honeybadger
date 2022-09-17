package evaluation

import (
	"testing"

	"github.com/notnil/chess"
	"github.com/stretchr/testify/assert"
)

func TestSimplified(t *testing.T) {
	tests := []struct {
		name string
		args string
		want int
	}{
		{name: "starting position", args: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", want: 0},
		{name: "endgame white", args: "8/8/8/5K1k/8/8/8/5R2 w - - 0 1", want: 480},
		{name: "endgame black", args: "7k/5K2/8/8/8/8/8/5R2 b - - 0 1", want: -440},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fen, err := chess.FEN(tt.args)
			assert.NoErrorf(t, err, "could not parse FEN %s", tt.args)

			game := chess.NewGame(fen)
			assert.Equal(t, tt.want, Simplified{}.Evaluate(game.Position()))
		})
	}
}

func TestMapSquareTableToWhite(t *testing.T) {
	want := [64]int{
		0, 0, 0, 0, 0, 0, 0, 0,
		5, 10, 10, -20, -20, 10, 10, 5,
		5, -5, -10, 0, 0, -10, -5, 5,
		0, 0, 0, 20, 20, 0, 0, 0,
		5, 5, 10, 25, 25, 10, 5, 5,
		10, 10, 20, 30, 30, 20, 10, 10,
		50, 50, 50, 50, 50, 50, 50, 50,
		0, 0, 0, 0, 0, 0, 0, 0,
	}
	assert.Equal(t, want, mapSquareTableToWhite(simplifiedHumanPawnTable, 0))
}

func TestMapSquareTableToBlack(t *testing.T) {
	want := [64]int{
		0, 0, 0, 0, 0, 0, 0, 0,
		50, 50, 50, 50, 50, 50, 50, 50,
		10, 10, 20, 30, 30, 20, 10, 10,
		5, 5, 10, 25, 25, 10, 5, 5,
		0, 0, 0, 20, 20, 0, 0, 0,
		5, -5, -10, 0, 0, -10, -5, 5,
		5, 10, 10, -20, -20, 10, 10, 5,
		0, 0, 0, 0, 0, 0, 0, 0,
	}
	assert.Equal(t, want, mapSquareTableToBlack(simplifiedHumanPawnTable, 0))
}
