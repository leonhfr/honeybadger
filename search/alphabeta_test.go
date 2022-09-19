package search

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/honeybadger/evaluation"
	"github.com/leonhfr/honeybadger/quiescence"
)

func TestAlphaBeta(t *testing.T) {
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
			want: want{Output{0, 1, -evaluation.Mate, 0, nil}, nil, nil},
		},
		{
			name: "mate in 1",
			args: args{"8/8/8/5K1k/8/8/8/5R2 w - - 0 1", 1},
			want: want{Output{1, 12, evaluation.Mate - 1, 0, nil}, []string{"f1h1"}, nil},
		},
		{
			name: "mate in 1",
			args: args{"r1b1kb1r/pppp1ppp/2n1pq2/8/3Pn2N/2P3P1/PP1NPP1P/R1BQKB1R b KQkq - 3 6", 1},
			want: want{Output{1, 3, evaluation.Mate - 1, 0, nil}, []string{"f6f2"}, nil},
		},
		{
			name: "mate in 2",
			args: args{"5rk1/pb2npp1/1pq4p/5p2/5B2/1B6/P2RQ1PP/2r1R2K b - - 0 1", 3},
			want: want{Output{3, 2995, evaluation.Mate - 3, 0, nil}, []string{"c6g2", "e2g2", "c1e1"}, nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o, err := alphaBeta(context.Background(), Input{
				Position:   position(tt.args.fen),
				Depth:      tt.args.depth,
				Evaluation: evaluation.Simplified{},
				Quiescence: quiescence.None{},
			}, -evaluation.Mate, evaluation.Mate)

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
