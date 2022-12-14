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

func TestPosition_MakeMove(t *testing.T) {
	for _, tt := range testPositions {
		t.Run(tt.moveUCI, func(t *testing.T) {
			pos := unsafeFEN(tt.preFEN)
			pos.MakeMove(tt.move)
			assert.Equal(t, tt.postFEN, pos.String())
		})
	}
}

func BenchmarkPosition_MakeMove(b *testing.B) {
	for _, bb := range testPositions {
		pos := unsafeFEN(bb.preFEN)
		b.Run(bb.moveUCI, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				meta, _ := pos.MakeMove(bb.move)
				pos.UnmakeMove(bb.move, meta)
			}
		})
	}
}

func TestPosition_UnmakeMove(t *testing.T) {
	for _, tt := range testPositions {
		t.Run(tt.moveUCI, func(t *testing.T) {
			pos := unsafeFEN(tt.preFEN)
			meta, _ := pos.MakeMove(tt.move)
			pos.UnmakeMove(tt.move, meta)
			assert.Equal(t, tt.preFEN, pos.String())
		})
	}
}

func BenchmarkPosition_PieceMap(b *testing.B) {
	for _, bb := range testPositions {
		pos := unsafeFEN(bb.preFEN)
		b.Run(bb.preFEN, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				pos.PieceMap(func(p Piece, sq Square) {})
			}
		})
	}
}
