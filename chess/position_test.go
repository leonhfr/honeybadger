package chess

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartingPosition(t *testing.T) {
	assert.Equal(t, startFEN, StartingPosition().String())
}

func TestFromFEN(t *testing.T) {
	tests := []struct {
		args string
		want error
	}{
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", nil},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPP/RNBQKBNR w KQkq - 0 1", errors.New("invalid fen rank field (PPPPPPP)")},
	}

	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			_, err := FromFEN(tt.args)
			assert.Equal(t, tt.want, err)
		})
	}
}

func TestPosition_Move(t *testing.T) {
	for _, tt := range testPositions {
		t.Run(tt.moveUCI, func(t *testing.T) {
			pos := unsafeFEN(tt.preFEN)
			got, _ := pos.Move(tt.move)
			assert.Equal(t, tt.postFEN, got.String())
		})
	}
}

func BenchmarkPosition_Move(b *testing.B) {
	for _, bb := range testPositions {
		b.Run(bb.moveUCI, func(b *testing.B) {
			pos := unsafeFEN(bb.preFEN)
			for n := 0; n < b.N; n++ {
				pos.Move(bb.move)
			}
		})
	}
}
