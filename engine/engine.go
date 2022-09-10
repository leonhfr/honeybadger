// Package engine implements the engine core.
package engine

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/leonhfr/honeybadger/search"
	"github.com/leonhfr/honeybadger/uci"
	"github.com/notnil/chess"
)

const (
	defaultMoveTime = 5 * time.Second
)

// Engine represents the engine object.
type Engine struct {
	name       string
	author     string
	game       *chess.Game
	mu         sync.Mutex
	stopSearch chan struct{}
	search     search.Interface
}

// New returns a new Engine.
func New(name, author string, options ...func(*Engine)) *Engine {
	e := &Engine{
		name:       name,
		author:     author,
		game:       chess.NewGame(),
		mu:         sync.Mutex{},
		stopSearch: make(chan struct{}),
	}

	for _, o := range availableOptions {
		fn := o.defaultFunc()
		fn(e)
	}

	for _, fn := range options {
		fn(e)
	}

	return e
}

// WithSearch sets the search strategy.
func WithSearch(si search.Interface) func(*Engine) {
	return func(e *Engine) {
		e.search = si
	}
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
	var options []uci.Option
	for _, option := range availableOptions {
		options = append(options, option.uci())
	}
	return options
}

// SetOption sets an option.
func (e *Engine) SetOption(name, value string) error {
	for _, option := range availableOptions {
		if option.String() == name {
			fn, err := option.optionFunc(value)
			if err != nil {
				return err
			}
			fn(e)
			return nil
		}
	}

	return errors.New("option name not found")
}

// SetPosition sets the position to the provided FEN.
func (e *Engine) SetPosition(fen string) error {
	fn, err := chess.FEN(fen)
	if err != nil {
		return err
	}
	fn(e.game)
	return nil
}

// Move plays the moves on the current position.
func (e *Engine) Move(moves ...*chess.Move) error {
	for _, move := range moves {
		if err := e.game.Move(move); err != nil {
			return err
		}
	}
	return nil
}

// ResetPosition resets the position to the starting one.
func (e *Engine) ResetPosition() {
	e.game = chess.NewGame()
}

// Search runs a search on the given input.
func (e *Engine) Search(input uci.Input) <-chan uci.Output {
	e.mu.Lock()
	start := time.Now()
	ctx, cancel := newContext(input, e.stopSearch)

	engineOutput := make(chan uci.Output)
	searchOutput := search.Run(ctx, search.Input{
		Position: e.game.Position(),
		Strategy: search.Random{},
	})

	go func() {
		defer e.mu.Unlock()
		defer cancel()
		defer close(engineOutput)

		for output := range searchOutput {
			engineOutput <- uci.Output{
				Time:  time.Since(start),
				Depth: output.Depth,
				Nodes: output.Nodes,
				Score: output.Score,
				PV:    output.PV,
			}
		}
	}()

	return engineOutput
}

// StopSearch aborts a search prematurely.
func (e *Engine) StopSearch() {
	select {
	case e.stopSearch <- struct{}{}:
	default:
	}
}

// Quit initiates a graceful shutdown.
func (e *Engine) Quit() {
	e.StopSearch()
	// prevents future searches and ensures all search routines have been shut down
	e.mu.Lock()
}

// newContext creates a new context from the input
func newContext(input uci.Input, stop <-chan struct{}) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	if !input.Infinite {
		timeout := moveTime(input)
		var unused context.CancelFunc
		ctx, unused = context.WithTimeout(ctx, timeout)
		_ = unused // pacify vet lostcancel check
	}

	go func() {
		select {
		case <-ctx.Done():
			return
		case <-stop:
			cancel()
		}
	}()

	return ctx, cancel
}

// moveTime determines how long the search should be
func moveTime(input uci.Input) time.Duration {
	if input.MoveTime > 0 {
		return input.MoveTime
	}

	return defaultMoveTime
}
