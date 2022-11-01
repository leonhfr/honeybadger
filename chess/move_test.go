package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkNewMove(b *testing.B) {
	for _, bb := range testPositions {
		p1, p2 := bb.move.P1(), bb.move.P2()
		s1, s2 := bb.move.S1(), bb.move.S2()
		promo := bb.move.Promo()
		ep := unsafeFEN(bb.preFEN).enPassant

		b.Run(bb.moveUCI, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				newMove(p1, p2, s1, s2, ep, promo)
			}
		})
	}
}

func TestFromUCI(t *testing.T) {
	for _, tt := range testPositions {
		t.Run(tt.moveUCI, func(t *testing.T) {
			move, _ := MoveFromUCI(unsafeFEN(tt.preFEN), tt.moveUCI)
			assert.Equal(t, tt.move, move)
			for _, tag := range tt.tags {
				assert.True(t, move.HasTag(tag))
			}
		})
	}
}

func TestFromUCI_Invalid(t *testing.T) {
	type (
		args struct {
			fen string
			uci string
		}
		want struct {
			move Move
			err  error
		}
	)

	tests := []struct {
		args args
		want want
	}{
		{
			args{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", "e2"},
			want{0, errInvalidMove},
		},
	}

	for _, tt := range tests {
		t.Run(tt.args.uci, func(t *testing.T) {
			move, err := MoveFromUCI(unsafeFEN(tt.args.fen), tt.args.uci)
			assert.Equal(t, tt.want.move, move)
			assert.Equal(t, tt.want.err, err)
		})
	}
}
