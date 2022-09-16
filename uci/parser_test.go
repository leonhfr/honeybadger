package uci

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"

	tests := []struct {
		name string
		args string
		want command
	}{
		{name: "uci", args: "uci", want: commandUCI{}},
		{name: "debug on", args: "debug on", want: commandDebug{on: true}},
		{name: "debug off", args: "debug off", want: commandDebug{on: false}},
		{name: "isready", args: "isready", want: commandIsReady{}},
		{name: "setoption", args: "setoption name NAME value VALUE", want: commandSetOption{name: "NAME", value: "VALUE"}},
		{name: "ucinewgame", args: "ucinewgame", want: commandUCINewGame{}},
		{name: "position", args: "position startpos", want: commandPosition{startPos: true}},
		{name: "position", args: "position fen " + fen, want: commandPosition{fen: fen}},
		{name: "position", args: "position startpos moves b1a3 b1c3", want: commandPosition{startPos: true, moves: []string{"b1a3", "b1c3"}}},
		{
			name: "go",
			args: "go infinite searchmoves b1a3 b1c3",
			want: commandGo{input: Input{
				SearchMoves: []string{"b1a3", "b1c3"},
				Infinite:    true,
			}},
		},
		{name: "stop", args: "stop", want: commandStop{}},
		{name: "quit", args: "quit", want: commandQuit{}},
		{name: "unknown", args: "foo bar", want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parse(strings.Fields(tt.args))
			assert.Equal(t, tt.want, got)
		})
	}
}
