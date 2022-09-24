package uci

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCommandUCI2(t *testing.T) {
	option := Option{Type: OptionBoolean, Name: "OPTION"}
	e := new(mockEngine)
	e.On("Info").Return("NAME", "AUTHOR")
	e.On("Options").Return([]Option{option})

	stdout := &strings.Builder{}
	respond := newResponder(stdout)

	commandUCI{}.run(context.Background(), e, respond)

	expected := concatenate([]response{
		responseID{"NAME", "AUTHOR"},
		option,
		responseUCIOK{},
	})
	e.AssertExpectations(t)
	assert.Equal(t, expected, stdout.String())
}

func TestCommandDebug(t *testing.T) {
	e := new(mockEngine)
	e.On("Debug", mock.Anything)

	stdout := &strings.Builder{}
	respond := newResponder(stdout)

	commandDebug{}.run(context.Background(), e, respond)

	e.AssertExpectations(t)
}

func TestCommandIsReady(t *testing.T) {
	tests := []struct {
		name string
		args error
		want []response
	}{
		{"no error", nil, []response{responseReadyOK{}}},
		{"error", errors.New("test"), []response{responseComment{"test"}, responseReadyOK{}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := new(mockEngine)
			e.On("Init").Return(tt.args)

			expected := concatenate(tt.want)
			stdout := newMockStdOut(len(expected))
			respond := newResponder(stdout)

			commandIsReady{}.run(context.Background(), e, respond)

			stdout.Wait()
			e.AssertExpectations(t)
			assert.Equal(t, expected, stdout.String())
		})
	}
}

func TestCommandSetOption(t *testing.T) {
	type args struct {
		cmd commandSetOption
		err error
	}

	tests := []struct {
		name string
		args args
		want []response
	}{
		{
			"valid option",
			args{commandSetOption{"NAME", "VALUE"}, nil},
			[]response{},
		},
		{
			"invalid option",
			args{commandSetOption{"NAME", "VALUE"}, errors.New("ERROR")},
			[]response{responseComment{"ERROR"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := new(mockEngine)
			e.On("SetOption", tt.args.cmd.name, tt.args.cmd.value).Return(tt.args.err)

			stdout := &strings.Builder{}
			respond := newResponder(stdout)

			tt.args.cmd.run(context.Background(), e, respond)

			e.AssertExpectations(t)
			assert.Equal(t, concatenate(tt.want), stdout.String())
		})
	}
}

func TestCommandUCINewGame(t *testing.T) {
	e := new(mockEngine)

	stdout := &strings.Builder{}
	respond := newResponder(stdout)

	commandUCINewGame{}.run(context.Background(), e, respond)

	e.AssertExpectations(t)
}

func TestCommandPosition_ResetPosition(t *testing.T) {
	type args struct {
		cmd     commandPosition
		errMove error
	}

	tests := []struct {
		name string
		args args
		want []response
	}{
		{
			"no error",
			args{commandPosition{startPos: true, moves: []string{"b1a3"}}, nil},
			[]response{},
		},
		{
			"error",
			args{commandPosition{startPos: true, moves: []string{"b1a3"}}, errors.New("ERROR")},
			[]response{responseComment{"ERROR"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := new(mockEngine)
			e.On("ResetPosition")
			e.On("Move", tt.args.cmd.moves).Return(tt.args.errMove)

			stdout := &strings.Builder{}
			respond := newResponder(stdout)

			tt.args.cmd.run(context.Background(), e, respond)

			e.AssertExpectations(t)
			assert.Equal(t, concatenate(tt.want), stdout.String())
		})
	}
}

func TestCommandPosition_SetPosition(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"

	type args struct {
		cmd     commandPosition
		errPos  error
		errMove error
	}

	tests := []struct {
		name string
		args args
		want []response
	}{
		{
			"no error",
			args{commandPosition{fen: fen, moves: []string{"b1a3"}}, nil, nil},
			[]response{},
		},
		{
			"error",
			args{commandPosition{fen: fen, moves: []string{"b1a3"}}, errors.New("ERROR POS"), errors.New("ERROR MOVE")},
			[]response{responseComment{"ERROR POS"}, responseComment{"ERROR MOVE"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := new(mockEngine)
			e.On("SetPosition", tt.args.cmd.fen).Return(tt.args.errPos)
			e.On("Move", tt.args.cmd.moves).Return(tt.args.errMove)

			stdout := &strings.Builder{}
			respond := newResponder(stdout)

			tt.args.cmd.run(context.Background(), e, respond)

			e.AssertExpectations(t)
			assert.Equal(t, concatenate(tt.want), stdout.String())
		})
	}
}

func TestCommandGo(t *testing.T) {
	type args struct {
		cmd     commandGo
		outputs []Output
		err     error
	}

	output1 := Output{Score: 1000, PV: []string{"b1a3", "d2d4"}}
	output2 := Output{Score: 2000, PV: []string{"d2d4"}}

	tests := []struct {
		name string
		args args
		want []response
	}{
		{
			"error",
			args{commandGo{Input{Depth: 3}}, []Output{}, errors.New("test")},
			[]response{responseComment{"test"}},
		},
		{
			"go",
			args{commandGo{Input{Depth: 3}}, []Output{output1, output2}, nil},
			[]response{output1, output2, responseBestMove{output2.PV[0]}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := new(mockEngine)
			oc := make(chan Output, len(tt.args.outputs))
			for _, o := range tt.args.outputs {
				oc <- o
			}
			close(oc)
			e.On("Search", mock.Anything, tt.args.cmd.input).Return(oc, tt.args.err)

			expected := concatenate(tt.want)
			stdout := newMockStdOut(len(expected))
			respond := newResponder(stdout)

			tt.args.cmd.run(context.Background(), e, respond)

			stdout.Wait()
			e.AssertExpectations(t)
			assert.Equal(t, expected, stdout.String())
		})
	}
}

func TestCommandStop(t *testing.T) {
	e := new(mockEngine)
	e.On("StopSearch")

	stdout := &strings.Builder{}
	respond := newResponder(stdout)

	commandStop{}.run(context.Background(), e, respond)

	e.AssertExpectations(t)
}

func TestCommandQuit(t *testing.T) {
	e := new(mockEngine)
	e.On("Quit")

	stdout := &strings.Builder{}
	respond := newResponder(stdout)

	commandQuit{}.run(context.Background(), e, respond)

	e.AssertExpectations(t)
}

func concatenate(responses []response) string {
	var s []string
	for _, r := range responses {
		s = append(s, r.String(), "\n")
	}
	return strings.Join(s, "")
}

type mockStdout struct {
	b   *strings.Builder
	wg  *sync.WaitGroup
	lim int
}

func newMockStdOut(lim int) *mockStdout {
	ms := &mockStdout{
		b:   &strings.Builder{},
		wg:  &sync.WaitGroup{},
		lim: lim,
	}
	ms.wg.Add(1)
	return ms
}

func (ms *mockStdout) Write(p []byte) (n int, err error) {
	n, err = ms.b.Write(p)
	if ms.b.Len() >= ms.lim {
		ms.wg.Done()
	}
	return n, err
}

func (ms *mockStdout) String() string {
	return ms.b.String()
}

func (ms *mockStdout) Wait() {
	ms.wg.Wait()
}
