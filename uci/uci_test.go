package uci

import (
	"github.com/stretchr/testify/mock"
)

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

func (m *mockEngine) Init() {
	m.Called()
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

func (m *mockEngine) Search(input Input) <-chan Output {
	args := m.Called(input)
	return args.Get(0).(chan Output)
}

func (m *mockEngine) StopSearch() {
	m.Called()
}

func (m *mockEngine) Quit() {
	m.Called()
}
