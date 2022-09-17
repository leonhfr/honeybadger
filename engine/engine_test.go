package engine

import (
	"errors"
	"testing"

	"github.com/notnil/chess"
	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/honeybadger/evaluation"
	"github.com/leonhfr/honeybadger/search"
	"github.com/leonhfr/honeybadger/uci"
)

func TestNew(t *testing.T) {
	e := New()
	assert.Equal(t, chess.StartingPosition().String(), e.game.Position().String())
	assert.Equal(t, chess.UCINotation{}, e.notation)
	assert.Equal(t, search.Negamax{}, e.search)
	assert.Equal(t, evaluation.Simplified{}, e.evaluation)
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
	assert.Equal(t, search.Capture{}, e.search)
}

func TestWithEvaluation(t *testing.T) {
	e := New(WithEvaluation(evaluation.Simplified{}))
	assert.Equal(t, evaluation.Simplified{}, e.evaluation)
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
			Default: "Negamax",
			Vars:    []string{"Capture", "Random", "Negamax"},
		},
		{
			Type:    uci.OptionEnum,
			Name:    "EvaluationStrategy",
			Default: "Simplified",
			Vars:    []string{"Values", "Simplified"},
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
			want{search.Random{}, errors.New("option value not found")},
		},
		{
			"option does not exist",
			args{"Whatever", "Whatever", search.Random{}},
			want{search.Random{}, errors.New("option name not found")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New(WithSearch(tt.args.search))
			err := e.SetOption(tt.args.name, tt.args.value)
			assert.Equal(t, tt.want.search, e.search)
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
