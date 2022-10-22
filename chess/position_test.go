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

func TestPosition_Update(t *testing.T) {
	for _, tt := range testPositions {
		t.Run(tt.move.String(), func(t *testing.T) {
			pos := unsafeFEN(tt.preFEN)
			got := pos.Move(tt.move)
			assert.Equal(t, tt.postFEN, got.String())
		})
	}
}
