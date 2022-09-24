package engine

import (
	"context"
	"errors"
	"testing"

	"github.com/notnil/chess"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/leonhfr/honeybadger/evaluation"
	"github.com/leonhfr/honeybadger/opening"
	"github.com/leonhfr/honeybadger/oracle"
	"github.com/leonhfr/honeybadger/quiescence"
	"github.com/leonhfr/honeybadger/search"
	"github.com/leonhfr/honeybadger/transposition"
	"github.com/leonhfr/honeybadger/uci"
)

func TestNew(t *testing.T) {
	e := New()
	assert.Equal(t, chess.StartingPosition().String(), e.game.Position().String())
	assert.Equal(t, chess.UCINotation{}, e.notation)
	assert.Equal(t, search.AlphaBeta{}, e.options.search)
	assert.Equal(t, evaluation.Simplified{}, e.options.evaluation)
	assert.Equal(t, oracle.Order{}, e.options.oracle)
	assert.Equal(t, quiescence.None{}, e.options.quiescence)
	assert.Equal(t, transposition.None{}, e.options.transposition)
	assert.Equal(t, opening.NewNone().String(), e.options.opening.String())
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

func TestInit(t *testing.T) {
	err := errors.New("test error")

	type want struct {
		err         error
		initialized bool
	}

	tests := []struct {
		name string
		args error
		want want
	}{
		{"no error", nil, want{nil, true}},
		{"error", err, want{err, false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := new(mockTransposition)
			e := New(WithTransposition(tr))
			tr.On("Init", e.options.hash).Return(tt.args).Times(1)

			err := e.Init()
			_ = e.Init() // test sync.Once
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.initialized, e.initialized)
		})
	}
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
			Name:    "OracleStrategy",
			Default: "Order",
			Vars:    []string{"None", "Order"},
		},
		{
			Type:    uci.OptionEnum,
			Name:    "QuiescenceStrategy",
			Default: "None",
			Vars:    []string{"None", "AlphaBeta"},
		},
		{
			Type:    uci.OptionEnum,
			Name:    "TranspositionStrategy",
			Default: "None",
			Vars:    []string{"None", "Ristretto"},
		},
		{
			Type:    uci.OptionEnum,
			Name:    "OpeningStrategy",
			Default: "None",
			Vars:    []string{"None"},
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
			initialized bool
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
			"engine not initialized",
			args{true, "SearchStrategy", "Capture", search.Random{}},
			want{search.Random{}, errSetOption},
		},
		{
			"option exists",
			args{false, "SearchStrategy", "Capture", search.Random{}},
			want{search.Capture{}, nil},
		},
		{
			"option does not exist",
			args{false, "SearchStrategy", "Whatever", search.Random{}},
			want{search.Random{}, errOptionValue},
		},
		{
			"option does not exist",
			args{false, "Whatever", "Whatever", search.Random{}},
			want{search.Random{}, errOptionName},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New(WithSearch(tt.args.search))
			e.initialized = tt.args.initialized
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

func TestSearch_Initialized(t *testing.T) {
	s := new(mockSearch)
	s.On("Search").Unset()
	e := New(WithSearch(s))
	_, err := e.Search(context.Background(), uci.Input{})
	assert.Equal(t, errSearch, err)
}

// mockSearch is a mock that implements search.Interface
type mockSearch struct {
	mock.Mock
}

func (m *mockSearch) String() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockSearch) Search(ctx context.Context, input search.Input, output chan<- *search.Output) {
	args := m.Called(ctx, input, output)
	_, len := args.Diff([]interface{}{})
	for i := 0; i < len; i++ {
		output <- args.Get(i).(*search.Output)
	}
}

// mockTransposition is a mock that implements transposition.Interface
type mockTransposition struct {
	mock.Mock
}

func (m *mockTransposition) String() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockTransposition) Init(size int) error {
	args := m.Called(size)
	return args.Error(0)
}

func (m *mockTransposition) Set(key *chess.Position, entry transposition.Entry) {
	m.Called(key, entry)
}

func (m *mockTransposition) Get(key *chess.Position) (transposition.Entry, bool) {
	args := m.Called(key)
	return args.Get(0).(transposition.Entry), args.Bool(1)
}

func (m *mockTransposition) Close() {
	m.Called()
}
