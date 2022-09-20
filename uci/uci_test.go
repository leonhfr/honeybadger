package uci

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRun(t *testing.T) {
	e := new(mockEngine)
	e.On("Info").Return("NAME", "AUTHOR")
	e.On("Options").Return([]Option{})
	e.On("Quit")

	r := strings.NewReader("uci\nfake command\nquit")
	w := &strings.Builder{}

	Run(context.Background(), e, r, w)

	e.AssertExpectations(t)
	assert.Equal(t, "id name NAME\nid author AUTHOR\nuciok\n", w.String())
}

// mockEngine is a mock that implements the Engine interface
type mockEngine struct {
	mock.Mock
}

func (m *mockEngine) Debug(on bool) {
	m.Called(on)
}

func (m *mockEngine) Info() (name, author string) {
	args := m.Called()
	return args.String(0), args.String(1)
}

func (m *mockEngine) Init() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockEngine) Options() []Option {
	args := m.Called()
	return args.Get(0).([]Option)
}

func (m *mockEngine) SetOption(name, value string) error {
	args := m.Called(name, value)
	return args.Error(0)
}

func (m *mockEngine) SetPosition(fen string) error {
	args := m.Called(fen)
	return args.Error(0)
}

func (m *mockEngine) Move(moves ...string) error {
	args := m.Called(moves)
	return args.Error(0)
}

func (m *mockEngine) ResetPosition() {
	m.Called()
}

func (m *mockEngine) Search(ctx context.Context, input Input) (<-chan Output, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(chan Output), args.Error(1)
}

func (m *mockEngine) StopSearch() {
	m.Called()
}

func (m *mockEngine) Quit() {
	m.Called()
}
