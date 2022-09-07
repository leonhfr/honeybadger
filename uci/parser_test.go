package uci

import (
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/leonhfr/honeybadger/engine"
	"github.com/notnil/chess"
)

func Test_Parse(t *testing.T) {
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
		want Command
	}{
		{name: "uci", args: "uci", want: CommandUCI{}},
		{name: "debug on", args: "debug on", want: CommandDebug{on: true}},
		{name: "debug off", args: "debug off", want: CommandDebug{on: false}},
		{name: "isready", args: "isready", want: CommandIsReady{}},
		{name: "setoption", args: "setoption name NAME value VALUE", want: CommandSetOption{name: "NAME", value: "VALUE"}},
		{name: "ucinewgame", args: "ucinewgame", want: CommandUCINewGame{}},
		{name: "position", args: "position startpos", want: CommandPosition{startPos: true}},
		{name: "position", args: "position fen " + fen, want: CommandPosition{fen: fen}},
		{name: "position", args: "position startpos moves b1a3 b1c3", want: CommandPosition{startPos: true, moves: []*chess.Move{move1, move2}}},
		{
			name: "go",
			args: "go infinite searchmoves b1a3 b1c3",
			want: CommandGo{input: engine.Input{
				SearchMoves: []*chess.Move{move1, move2},
				Infinite:    true,
			}},
		},
		{name: "stop", args: "stop", want: CommandStop{}},
		{name: "quit", args: "quit", want: CommandQuit{}},
		{name: "unknown", args: "foo bar", want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse(strings.Fields(tt.args)); !cmp.Equal(got, tt.want, cmpOptions...) {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
