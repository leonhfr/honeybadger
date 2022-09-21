// Package uci implements the Universal Chess Interface.
package uci

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"
)

// Engine is the interface implemented by objects that can be used as UCI engines.
type Engine interface {
	Init() error                                                    // Init sets everything up.
	Quit()                                                          // Quit initiates a graceful shutdown.
	Debug(on bool)                                                  // Debug sets the debug option.
	Info() (name, author string)                                    // Info returns the engine's info.
	Options() []Option                                              // Options lists the available options.
	SetOption(name, value string) error                             // SetOption sets an option.
	SetPosition(fen string) error                                   // SetPosition sets the position to the provided FEN.
	ResetPosition()                                                 // ResetPosition resets the position to the starting one.
	Move(moves ...string) error                                     // Move plays the moves on the current position.
	Search(ctx context.Context, input Input) (<-chan Output, error) // Search runs a search on the given input.
	StopSearch()                                                    // StopSearch aborts a search prematurely.
}

// Run runs the program in UCI mode.
//
// Run parses command from the reader, executes them with the provided engine
// and writes the responses on the writer.
func Run(ctx context.Context, e Engine, r io.Reader, w io.Writer) {
	responses := make(chan response)
	wg := sync.WaitGroup{}
	wg.Add(1)

	logger := log.New(w, "", 0)
	go func() {
		defer wg.Done()
		for response := range responses {
			logger.Println(response)
		}
	}()

	for scanner := bufio.NewScanner(r); scanner.Scan(); {
		c := parse(strings.Fields(scanner.Text()))
		if c == nil {
			continue
		}
		c.run(ctx, e, responses)
		if _, ok := c.(commandQuit); ok {
			break
		}
	}

	close(responses)
	wg.Wait()
}

// Logger returns a logger that is able to log UCI-compliant output.
func Logger(w io.Writer) *log.Logger {
	return log.New(w, "info string ", 0)
}

// Input is what the engine needs to run a search.
type Input struct {
	WhiteTime      time.Duration // White has <x> ms left on the clock.
	BlackTime      time.Duration // Black has <x> ms left on the clock.
	WhiteIncrement time.Duration // White increment per move in ms if <x> > 0.
	BlackIncrement time.Duration // Black increment per move in ms if <x> > 0.
	MovesToGo      int           // Number of moves until the next time control.
	SearchMoves    []string      // Restrict search to those moves only.
	Depth          int           // Search <x> plies only.
	Nodes          int           // Search <x> nodes only.
	MoveTime       time.Duration // Search exactly <x> ms.
	Infinite       bool          // Search until the stop command. Do not exit before.
}

func (i Input) String() string {
	var res []string
	if i.Depth > 0 {
		res = append(res, fmt.Sprintf("depth %v", i.Depth))
	}
	if i.MoveTime > 0 {
		res = append(res, fmt.Sprintf("movetime %v", i.MoveTime.Milliseconds()))
	}
	if i.Infinite {
		res = append(res, "infinite")
	}
	if len(i.SearchMoves) > 0 {
		res = append(res, fmt.Sprintf("searchmoves %s", strings.Join(i.SearchMoves, " ")))
	}
	return fmt.Sprintf("go %v", strings.Join(res, " "))
}

// Output holds a search result.
type Output struct {
	Time  time.Duration // Time searched in ms.
	Depth int           // Search depth in plies.
	Nodes int           // Number of nodes searched.
	Score int           // Score from the engine's point of view in centipawns.
	Mate  int           // Number of moves before mate.
	PV    []string      // Principal variation, best line found.
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
