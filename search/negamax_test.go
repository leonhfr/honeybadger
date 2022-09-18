package search

import (
	"context"
	"testing"

	"github.com/notnil/chess"
	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/honeybadger/evaluation"
)

func TestNegamax(t *testing.T) {
	type (
		args struct {
			fen   string
			depth int
		}
		want struct {
			output Output
			moves  []string
			err    error
		}
	)

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "checkmate",
			args: args{"8/8/8/5K1k/8/8/8/7R b - - 0 1", 1},
			want: want{Output{0, 1, -mateScore, 0, nil}, nil, nil},
		},
		{
			name: "mate in 1",
			args: args{"8/8/8/5K1k/8/8/8/5R2 w - - 0 1", 1},
			want: want{Output{1, 15, mateScore - 1, 0, nil}, []string{"f1h1"}, nil},
		},
		{
			name: "mate in 1",
			args: args{"r1b1kb1r/pppp1ppp/2n1pq2/8/3Pn2N/2P3P1/PP1NPP1P/R1BQKB1R b KQkq - 3 6", 1},
			want: want{Output{1, 46, mateScore - 1, 0, nil}, []string{"f6f2"}, nil},
		},
		{
			name: "mate in 2",
			args: args{"5rk1/pb2npp1/1pq4p/5p2/5B2/1B6/P2RQ1PP/2r1R2K b - - 0 1", 3},
			want: want{Output{3, 90094, 9223372036854775804, 0, nil}, []string{"c6g2", "e2g2", "c1e1"}, nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o, err := negamax(context.Background(), Input{
				Position:   position(tt.args.fen),
				Depth:      tt.args.depth,
				Evaluation: evaluation.Simplified{},
			})

			output := *o
			pv := output.PV
			output.PV = nil

			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.output, output)
			assert.Len(t, pv, len(tt.want.moves))
			if len(pv) == len(tt.want.moves) {
				for i, move := range pv {
					assert.Equal(t, tt.want.moves[i], move.String())
				}
			}
		})
	}
}

func TestUpdateScore(t *testing.T) {
	tests := []struct {
		name string
		args int
		want int
	}{
		{"", mateScore, mateScore - 1},
		{"", mateScore - 1, mateScore - 2},
		{"", mateScore - 2, mateScore - 3},
		{"", -mateScore, -mateScore + 1},
		{"", -mateScore + 1, -mateScore + 2},
		{"", -mateScore + 2, -mateScore + 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, updateScore(tt.args))
		})
	}
}

func TestMatePlies(t *testing.T) {
	tests := []struct {
		name string
		args int
		want int
	}{
		{"normal move", 0, 0},
		{"normal move", 100, 0},
		{"normal move", -100, 0},
		{"mate", mateScore, 0},
		{"mate", -mateScore, 0},
		{"mate in 1", mateScore - 1, 1},
		{"mate in 1", -mateScore + 1, -1},
		{"mate in 2", mateScore - 2, 1},
		{"mate in 2", -mateScore + 2, -1},
		{"mate in 2", mateScore - 3, 2},
		{"mate in 2", -mateScore + 3, -2},
		{"mate in 2", mateScore - 4, 2},
		{"mate in 2", -mateScore + 4, -2},
		{"mate in 3", mateScore - 5, 3},
		{"mate in 3", -mateScore + 5, -3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, matePlies(tt.args))
		})
	}
}

func TestTerminalNode(t *testing.T) {
	type want struct {
		isNil  bool
		output Output
	}

	tests := []struct {
		name string
		args string // fen
		want want
	}{
		{
			name: "in game",
			args: "8/8/8/5K1k/8/8/8/7R w - - 0 1",
			want: want{true, Output{}},
		},
		{
			name: "checkmate",
			args: "8/8/8/5K1k/8/8/8/7R b - - 0 1",
			want: want{false, Output{Nodes: 1, Score: -mateScore}},
		},
		// TODO: stalemate (FEN: 7k/5K2/8/5R2/8/8/8/8 b - - 0 1) github.com/notnil/chess doesn't decode check and valid moves from FEN
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := terminalNode(position(tt.args))
			if tt.want.isNil {
				assert.Nil(t, output)
			} else {
				assert.NotNil(t, output)
				assert.Equal(t, tt.want.output, *output)
			}
		})
	}
}

func position(fen string) *chess.Position {
	fn, _ := chess.FEN(fen)
	game := chess.NewGame(fn)
	game.Outcome()
	return game.Position()
}
