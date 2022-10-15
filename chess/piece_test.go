package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPiece_Color(t *testing.T) {
	tests := []struct {
		name string
		args Piece
		want Color
	}{
		{"WhitePawn", WhitePawn, White},
		{"WhiteKnight", WhiteKnight, White},
		{"WhiteBishop", WhiteBishop, White},
		{"WhiteRook", WhiteRook, White},
		{"WhiteQueen", WhiteQueen, White},
		{"WhiteKing", WhiteKing, White},
		{"BlackPawn", BlackPawn, Black},
		{"BlackKnight", BlackKnight, Black},
		{"BlackBishop", BlackBishop, Black},
		{"BlackRook", BlackRook, Black},
		{"BlackQueen", BlackQueen, Black},
		{"BlackKing", BlackKing, Black},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.Color())
		})
	}
}

func TestPiece_Type(t *testing.T) {
	tests := []struct {
		name string
		args Piece
		want PieceType
	}{
		{"WhitePawn", WhitePawn, Pawn},
		{"WhiteKnight", WhiteKnight, Knight},
		{"WhiteBishop", WhiteBishop, Bishop},
		{"WhiteRook", WhiteRook, Rook},
		{"WhiteQueen", WhiteQueen, Queen},
		{"WhiteKing", WhiteKing, King},
		{"BlackPawn", BlackPawn, Pawn},
		{"BlackKnight", BlackKnight, Knight},
		{"BlackBishop", BlackBishop, Bishop},
		{"BlackRook", BlackRook, Rook},
		{"BlackQueen", BlackQueen, Queen},
		{"BlackKing", BlackKing, King},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.Type())
		})
	}
}
