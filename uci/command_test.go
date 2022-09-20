package uci

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCommandUCI(t *testing.T) {
	option := Option{Type: OptionBoolean, Name: "OPTION"}
	e := new(mockEngine)
	e.On("Info").Return("NAME", "AUTHOR")
	e.On("Options").Return([]Option{option})

	rc := make(chan response)
	wg := assertResponses(t, e, rc, []response{
		responseID{"NAME", "AUTHOR"},
		responseOption{option},
		responseUCIOK{},
	})

	commandUCI{}.run(e, rc)
	close(rc)

	e.AssertExpectations(t)
	wg.Wait()
}

func TestCommandDebug(t *testing.T) {
	e := new(mockEngine)
	e.On("Debug", mock.Anything)

	rc := make(chan response)

	commandDebug{}.run(e, rc)
	close(rc)

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

			rc := make(chan response)
			wg := &sync.WaitGroup{}
			wg.Add(1)

			// testing asynchronously
			go func() {
				defer wg.Done()
				defer close(rc)
				var r []response

				for i := 0; i < len(tt.want); i++ {
					r = append(r, <-rc)
				}

				select {
				case empty := <-rc:
					assert.Fail(t, "channel should be empty", empty)
				default:
				}

				e.AssertExpectations(t)
				assert.Equal(t, tt.want, r)
			}()

			commandIsReady{}.run(e, rc)

			wg.Wait()
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

			rc := make(chan response)
			wg := assertResponses(t, e, rc, tt.want)

			tt.args.cmd.run(e, rc)
			close(rc)

			e.AssertExpectations(t)
			wg.Wait()
		})
	}
}

func TestCommandUCINewGame(t *testing.T) {
	e := new(mockEngine)

	rc := make(chan response)

	commandUCINewGame{}.run(e, rc)
	close(rc)

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

			rc := make(chan response)
			wg := assertResponses(t, e, rc, tt.want)

			tt.args.cmd.run(e, rc)
			close(rc)

			e.AssertExpectations(t)
			wg.Wait()
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

			rc := make(chan response)
			wg := assertResponses(t, e, rc, tt.want)

			tt.args.cmd.run(e, rc)
			close(rc)

			e.AssertExpectations(t)
			wg.Wait()
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
			[]response{responseInfo{output1}, responseInfo{output2}, responseBestMove{output2.PV[0]}},
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
			e.On("Search", tt.args.cmd.input).Return(oc, tt.args.err)

			rc := make(chan response)
			wg := &sync.WaitGroup{}
			wg.Add(1)

			// testing asynchronously
			go func() {
				defer wg.Done()
				defer close(rc)
				var responses []response
				for i := 0; i < len(tt.want); i++ {
					responses = append(responses, <-rc)
				}
				e.AssertExpectations(t)
				assert.Equal(t, tt.want, responses)
			}()

			tt.args.cmd.run(e, rc)

			wg.Wait()
		})
	}
}

func TestCommandStop(t *testing.T) {
	e := new(mockEngine)
	e.On("StopSearch")

	rc := make(chan response)

	commandStop{}.run(e, rc)
	close(rc)

	e.AssertExpectations(t)
}

func TestCommandQuit(t *testing.T) {
	e := new(mockEngine)
	e.On("Quit")

	rc := make(chan response)

	commandQuit{}.run(e, rc)
	close(rc)

	e.AssertExpectations(t)
}

func assertResponses(t *testing.T, e *mockEngine, rc chan response, expected []response) *sync.WaitGroup {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		responses := []response{}
		for r := range rc {
			responses = append(responses, r)
		}
		assert.Equal(t, expected, responses)
	}()

	return wg
}
