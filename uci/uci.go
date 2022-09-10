// Package uci implements the Universal Chess Interface.
package uci

import (
	"bufio"
	"io"
	"log"
	"strings"
	"time"

	"github.com/notnil/chess"
)

// Engine is the interface implemented by objects that can be used as UCI engines.
type Engine interface {
	Init()
	Quit()
	Debug(on bool)
	Info() (name, author string)
	Options() []Option
	SetOption(name, value string) error
	SetPosition(fen string) error
	ResetPosition()
	Move(moves ...*chess.Move) error
	Search(input Input) <-chan Output
	StopSearch()
}

// Run runs the program in UCI mode.
//
// Run parses command from the reader, executes them with the provided engine
// and writes the responses on the writer.
func Run(e Engine, r io.Reader, w io.Writer) {
	responses := make(chan response)
	defer close(responses)

	logger := log.New(w, "", 0)
	go func() {
		for response := range responses {
			logger.Println(response)
		}
	}()

	for scanner := bufio.NewScanner(r); scanner.Scan(); {
		c := parse(strings.Fields(scanner.Text()))
		c.run(e, responses)
		if _, ok := c.(commandQuit); ok {
			break
		}
	}
}

// Input is what the engine needs to run a search.
type Input struct {
	WhiteTime      time.Duration // White has <x> ms left on the clock.
	BlackTime      time.Duration // Black has <x> ms left on the clock.
	WhiteIncrement time.Duration // White increment per move in ms if <x> > 0.
	BlackIncrement time.Duration // Black increment per move in ms if <x> > 0.
	MovesToGo      int           // Number of moves until the next time control.
	SearchMoves    []*chess.Move // Restrict search to those moves only.
	Depth          int           // Search <x> plies only.
	Nodes          int           // Search <x> nodes only.
	MoveTime       time.Duration // Search exactly <x> ms.
	Infinite       bool          // Search until the stop command. Do not exit before.
}

// Output holds a search result.
type Output struct {
	Time  time.Duration // Time searched in ms.
	Depth int           // Search depth in plies.
	Nodes int           // Number of nodes searched.
	Score int           // Score from the engine's point of view in centipawns.
	Mate  int           // Number of moves before mate.
	PV    []*chess.Move // Principal variation, best line found.
}

// OptionType represents an option's type.
type OptionType int

const (
	OptionBoolean OptionType = iota // OptionBoolean represents a boolean option.
	OptionInteger                   // OptionInteger represents an integer option.
	OptionEnum                      // OptionEnum represents an enum option.
)

// Option represents an available option.
type Option struct {
	Type    OptionType
	Name    string
	Default string
	Min     string
	Max     string
	Vars    []string
}
