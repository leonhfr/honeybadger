// Package engine implements the engine core.
package engine

import (
	"time"

	"github.com/notnil/chess"
)

// Engine represents the engine object.
type Engine struct {
	info Info
}

// Info holds the engine information.
type Info struct {
	Name    string
	Version string
	Author  string
}

// New returns a new Engine.
func New(info Info) *Engine {
	e := &Engine{
		info: info,
	}

	return e
}

// Debug sets the debug option.
func (e *Engine) Debug(on bool) {}

// Info returns the engine's info.
func (e *Engine) Info() Info {
	return e.info
}

// Initialize sets everything up.
func (e *Engine) Initialize() {}

// SetOption sets an option.
func (e *Engine) SetOption(name, value string) error {
	return nil
}

// SetFEN sets the position to the provided FEN.
func (e *Engine) SetFEN(fen string) error {
	return nil
}

// Move plays the moves on the current position.
func (e *Engine) Move(moves ...*chess.Move) error {
	return nil
}

// ResetPosition resets the position to the starting one.
func (e *Engine) ResetPosition() {
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

// Output is what the engine sends back,
// it represents a search result.
type Output struct {
	Done  bool          // Search is done. Best result has been found.
	Time  time.Duration // Time searched in ms.
	Depth int           // Search depth in plies.
	Nodes int           // Number of nodes searched.
	Score int           // Score from the engine's point of view in centipawns.
	Mate  int           // Number of moves before mate.
	PV    []*chess.Move // Principal variation, best line found.
}

// Search runs a search on the given input.
func (e *Engine) Search(input Input) <-chan Output {
	engineOutput := make(chan Output)
	defer close(engineOutput)
	return engineOutput
}

// StopSearch stops a search prematurely.
func (e *Engine) StopSearch() {
}

// Quit initiates a graceful shutdown.
func (e *Engine) Quit() {
}
