package uci

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/notnil/chess"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
	moves := chess.StartingPosition().ValidMoves()
	move1 := moves[0]
	move2 := moves[1]

	cmpOptions := []cmp.Option{
		cmp.Exporter(func(t reflect.Type) bool { return true }),
		cmp.Transformer("", func(move *chess.Move) string {
			return chess.UCINotation{}.Encode(nil, move)
		}),
	}

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
		{name: "position", args: "position startpos moves b1a3 b1c3", want: commandPosition{startPos: true, moves: []*chess.Move{move1, move2}}},
		{
			name: "go",
			args: "go infinite searchmoves b1a3 b1c3",
			want: commandGo{input: Input{
				SearchMoves: []*chess.Move{move1, move2},
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
			assert.True(
				t,
				cmp.Equal(got, tt.want, cmpOptions...),
				fmt.Sprintf("expected %v, got %v", tt.want, got),
			)
		})
	}
}
