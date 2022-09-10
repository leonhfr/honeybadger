// Package engine implements the engine core.
package engine

import (
	"github.com/leonhfr/honeybadger/uci"
	"github.com/notnil/chess"
)

// Engine represents the engine object.
type Engine struct {
	name   string
	author string
}

// New returns a new Engine.
func New(name, author string) *Engine {
	e := &Engine{
		name:   name,
		author: author,
	}

	return e
}

// Debug sets the debug option.
func (e *Engine) Debug(on bool) {}

// Info returns the engine's info.
func (e *Engine) Info() (name, author string) {
	return e.name, e.author
}

// Init sets everything up.
func (e *Engine) Init() {}

// Options lists the available options.
func (e *Engine) Options() []uci.Option {
	return nil
}

// SetOption sets an option.
func (e *Engine) SetOption(name, value string) error {
	return nil
}

// SetPosition sets the position to the provided FEN.
func (e *Engine) SetPosition(fen string) error {
	return nil
}

// Move plays the moves on the current position.
func (e *Engine) Move(moves ...*chess.Move) error {
	return nil
}

// ResetPosition resets the position to the starting one.
func (e *Engine) ResetPosition() {
}

// Search runs a search on the given input.
func (e *Engine) Search(input uci.Input) <-chan uci.Output {
	engineOutput := make(chan uci.Output)
	defer close(engineOutput)
	return engineOutput
}

// StopSearch stops a search prematurely.
func (e *Engine) StopSearch() {
}

// Quit initiates a graceful shutdown.
func (e *Engine) Quit() {
}
