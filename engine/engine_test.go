package engine

import (
	"testing"

	"github.com/notnil/chess"
	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/honeybadger/evaluation"
	"github.com/leonhfr/honeybadger/quiescence"
	"github.com/leonhfr/honeybadger/search"
	"github.com/leonhfr/honeybadger/uci"
)

func TestNew(t *testing.T) {
	e := New()
	assert.Equal(t, chess.StartingPosition().String(), e.game.Position().String())
	assert.Equal(t, chess.UCINotation{}, e.notation)
	assert.Equal(t, search.AlphaBeta{}, e.options.search)
	assert.Equal(t, evaluation.Simplified{}, e.options.evaluation)
	assert.Equal(t, quiescence.None{}, e.options.quiescence)
	assert.Equal(t, 32, e.options.hash)
}

func TestWithName(t *testing.T) {
	e := New(WithName("NAME"))
	assert.Equal(t, "NAME", e.name)
}

func TestWithAuthor(t *testing.T) {
	e := New(WithAuthor("AUTHOR"))
	assert.Equal(t, "AUTHOR", e.author)
}

func TestWithSearch(t *testing.T) {
	e := New(WithSearch(search.Capture{}))
	assert.Equal(t, search.Capture{}, e.options.search)
}

func TestWithEvaluation(t *testing.T) {
	e := New(WithEvaluation(evaluation.Simplified{}))
	assert.Equal(t, evaluation.Simplified{}, e.options.evaluation)
}

func TestWithQuiescence(t *testing.T) {
	e := New(WithQuiescence(quiescence.AlphaBeta{}))
	assert.Equal(t, quiescence.AlphaBeta{}, e.options.quiescence)
}

func TestInfo(t *testing.T) {
	e := New(WithName("NAME"), WithAuthor("AUTHOR"))
	name, author := e.Info()
	assert.Equal(t, "NAME", name)
	assert.Equal(t, "AUTHOR", author)
}

func TestOptions(t *testing.T) {
	e := New()
	options := e.Options()
	assert.Equal(t, []uci.Option{
		{
			Type:    uci.OptionEnum,
			Name:    "SearchStrategy",
			Default: "AlphaBeta",
			Vars:    []string{"Random", "Capture", "Negamax", "AlphaBeta"},
		},
		{
			Type:    uci.OptionEnum,
			Name:    "EvaluationStrategy",
			Default: "Simplified",
			Vars:    []string{"Values", "Simplified"},
		},
		{
			Type:    uci.OptionEnum,
			Name:    "QuiescenceStrategy",
			Default: "None",
			Vars:    []string{"None", "AlphaBeta"},
		},
		{
			Type:    uci.OptionInteger,
			Name:    "Hash",
			Default: "32",
			Min:     "1",
			Max:     "1024",
		},
	}, options)
}

func TestSetOption(t *testing.T) {
	type (
		args struct {
			name, value string
			search      search.Interface
		}
		want struct {
			search search.Interface
			err    error
		}
	)

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"option exists",
			args{"SearchStrategy", "Capture", search.Random{}},
			want{search.Capture{}, nil},
		},
		{
			"option does not exist",
			args{"SearchStrategy", "Whatever", search.Random{}},
			want{search.Random{}, errOptionValue},
		},
		{
			"option does not exist",
			args{"Whatever", "Whatever", search.Random{}},
			want{search.Random{}, errOptionName},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New(WithSearch(tt.args.search))
			err := e.SetOption(tt.args.name, tt.args.value)
			assert.Equal(t, tt.want.search, e.options.search)
			assert.Equal(t, tt.want.err, err)
		})
	}
}

func TestSetPositionValid(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
	e := New()
	err := e.SetPosition(fen)
	if assert.NoError(t, err) {
		assert.Equal(t, fen, e.game.Position().String())
	}
}

func TestSetPositionInvalid(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/BROKEN_ b KQkq e3 0 1"
	e := New()
	err := e.SetPosition(fen)
	assert.Error(t, err)
}

func TestMoveValid(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
	e := New()
	err := e.Move("e2e4")
	if assert.NoError(t, err) {
		assert.Equal(t, fen, e.game.Position().String())
	}
}

func TestMoveInvalidDecode(t *testing.T) {
	e := New()
	err := e.Move("e2e")
	assert.Error(t, err)
}

func TestMoveInvalidMove(t *testing.T) {
	e := New()
	err := e.Move("e2e5")
	assert.Error(t, err)
}

func TestResetPosition(t *testing.T) {
	e := New()
	e.ResetPosition()
	assert.Equal(t, chess.StartingPosition().String(), e.game.Position().String())
}
