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
	wg := &sync.WaitGroup{}

	assertResponses(t, wg, e, rc, []response{
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
	e := new(mockEngine)
	e.On("Init")

	rc := make(chan response)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer close(rc)
		r := <-rc
		// testing asynchronously
		e.AssertExpectations(t)
		assert.Equal(t, responseReadyOK{}, r)
	}()

	commandIsReady{}.run(e, rc)

	wg.Wait()
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
			wg := &sync.WaitGroup{}

			assertResponses(t, wg, e, rc, tt.want)

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

// TODO: position

// TODO: go

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

func assertResponses(t *testing.T, wg *sync.WaitGroup, e *mockEngine, rc chan response, expected []response) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		responses := []response{}
		for r := range rc {
			responses = append(responses, r)
		}
		assert.Equal(t, expected, responses)
	}()
}
