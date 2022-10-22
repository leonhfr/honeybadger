package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromUCI(t *testing.T) {
	type (
		args struct {
			fen string
			uci string
		}
		want struct {
			move Move
			tags []MoveTag
			err  error
		}
	)

	tests := []struct {
		args args
		want want
	}{
		{
			args{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", "e2e4"},
			want{newMove(WhitePawn, NoPiece, E2, E4, NoSquare, NoPiece), []MoveTag{Quiet}, nil},
		},
		{
			args{"r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1", "e1g1"},
			want{newMove(WhiteKing, NoPiece, E1, G1, NoSquare, NoPiece), []MoveTag{KingSideCastle}, nil},
		},
		{
			args{"2r3k1/1q1nbppp/r3p3/3pP3/pPpP4/P1Q2N2/2RN1PPP/2R4K b - b3 0 23", "a4b3"},
			want{newMove(BlackPawn, NoPiece, A4, B3, B3, NoPiece), []MoveTag{EnPassant, Capture}, nil},
		},
		{
			args{"r2qk2r/pp1n1ppp/2pbpn2/3p4/2PP4/1PNQPN2/P4PPP/R1B1K2R w KQkq - 1 9", "e1g1"},
			want{newMove(WhiteKing, NoPiece, E1, G1, NoSquare, NoPiece), []MoveTag{KingSideCastle}, nil},
		},
		{
			args{"r1bqkbnr/ppp1pppp/2n5/3p4/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3", "e4d5"},
			want{newMove(WhitePawn, BlackPawn, E4, D5, NoSquare, NoPiece), []MoveTag{Capture}, nil},
		},
		{
			args{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", "e2"},
			want{0, []MoveTag{}, errInvalidMove},
		},
	}

	for _, tt := range tests {
		t.Run(tt.args.uci, func(t *testing.T) {
			move, err := MoveFromUCI(unsafeFEN(tt.args.fen), tt.args.uci)
			assert.Equal(t, tt.want.move, move)
			for _, tag := range tt.want.tags {
				assert.True(t, move.HasTag(tag))
			}
			assert.Equal(t, tt.want.err, err)
		})
	}
}

func unsafeFEN(fen string) *Position {
	p, err := FromFEN(fen)
	if err != nil {
		panic(err)
	}
	return p
}
