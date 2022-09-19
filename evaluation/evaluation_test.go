package evaluation

import (
	"testing"

	"github.com/notnil/chess"
	"github.com/stretchr/testify/assert"
)

func TestTerminal(t *testing.T) {
	type want struct {
		score    int
		terminal bool
	}

	tests := []struct {
		name string
		args string // fen
		want want
	}{
		{
			name: "in game",
			args: "8/8/8/5K1k/8/8/8/7R w - - 0 1",
			want: want{0, false},
		},
		{
			name: "checkmate",
			args: "8/8/8/5K1k/8/8/8/7R b - - 0 1",
			want: want{-Mate, true},
		},
		// TODO: stalemate (FEN: 7k/5K2/8/5R2/8/8/8/8 b - - 0 1) github.com/notnil/chess doesn't decode check and valid moves from FEN
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, terminal := Terminal(position(tt.args))
			assert.Equal(t, tt.want.score, score)
			assert.Equal(t, tt.want.terminal, terminal)
		})
	}
}

func position(fen string) *chess.Position {
	fn, _ := chess.FEN(fen)
	game := chess.NewGame(fn)
	game.Outcome()
	return game.Position()
}
