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
			move *Move
			err  error
		}
	)

	tests := []struct {
		args args
		want want
	}{
		{
			args{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", "e2e4"},
			want{&Move{s1: E2, s2: E4, promo: NoPieceType}, nil},
		},
		{
			args{"r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1", "e1g1"},
			want{&Move{s1: E1, s2: G1, promo: NoPieceType, tags: KingSideCastle}, nil},
		},
		{
			args{"2r3k1/1q1nbppp/r3p3/3pP3/pPpP4/P1Q2N2/2RN1PPP/2R4K b - b3 0 23", "a4b3"},
			want{&Move{s1: A4, s2: B3, promo: NoPieceType, tags: Capture | EnPassant}, nil},
		},
		{
			args{"r2qk2r/pp1n1ppp/2pbpn2/3p4/2PP4/1PNQPN2/P4PPP/R1B1K2R w KQkq - 1 9", "e1g1"},
			want{&Move{s1: E1, s2: G1, promo: NoPieceType, tags: KingSideCastle}, nil},
		},
		{
			args{"r1bqkbnr/ppp1pppp/2n5/3p4/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3", "e4d5"},
			want{&Move{s1: E4, s2: D5, promo: NoPieceType, tags: Capture}, nil},
		},
		{
			args{"1k6/8/8/8/8/8/4R3/4K3 w - - 0 1", "e2b2"},
			want{&Move{s1: E2, s2: B2, promo: NoPieceType, tags: Check}, nil},
		},
		{
			args{"1k2q3/8/8/8/8/8/4R3/4K3 w - - 0 1", "e2b2"},
			want{nil, errIllegalMove},
		},
		{
			args{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", "e2"},
			want{nil, errInvalidMove},
		},
	}

	for _, tt := range tests {
		t.Run(tt.args.uci, func(t *testing.T) {
			move, err := FromUCI(unsafeFEN(tt.args.fen), tt.args.uci)
			if tt.want.move != nil {
				assert.Equal(t, *tt.want.move, *move)
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
