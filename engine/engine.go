// Package engine implements the engine core.
package engine

import (
	"context"
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"github.com/notnil/chess"

	"github.com/leonhfr/honeybadger/evaluation"
	"github.com/leonhfr/honeybadger/opening"
	"github.com/leonhfr/honeybadger/oracle"
	"github.com/leonhfr/honeybadger/quiescence"
	"github.com/leonhfr/honeybadger/search"
	"github.com/leonhfr/honeybadger/transposition"
	"github.com/leonhfr/honeybadger/uci"
)

const (
	defaultMoveTime = 5 * time.Second
)

var (
	errOptionName = errors.New("option name not found")
	errSetOption  = errors.New("cannot set option after engine has been initialized")
	errSearch     = errors.New("cannot run a search before engine has been initialized")
)

// Engine represents the engine object.
type Engine struct {
	name        string
	author      string
	debug       bool
	logger      *log.Logger
	game        *chess.Game
	notation    chess.Notation
	mu          sync.Mutex
	once        sync.Once
	initialized bool
	stopSearch  chan struct{}
	options     engineOptions
}

type engineOptions struct {
	search        search.Interface        // Search strategy.
	evaluation    evaluation.Interface    // Evaluation strategy.
	oracle        oracle.Interface        // Oracle strategy.
	quiescence    quiescence.Interface    // Quiescence strategy.
	transposition transposition.Interface // Transposition strategy.
	opening       opening.Interface       // Opening strategy.
	hash          int                     // Size of the transposition hash table in MB.
}

// New returns a new Engine.
//
// Defaults to using the UCI chess notation and a
// logger that outputs to stdout without any prefix.
func New(options ...func(*Engine)) *Engine {
	e := &Engine{
		logger:     log.New(os.Stdout, "", 0),
		game:       chess.NewGame(),
		notation:   chess.UCINotation{},
		mu:         sync.Mutex{},
		stopSearch: make(chan struct{}),
		options:    engineOptions{},
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

// WithName sets the name of the engine.
func WithName(name string) func(*Engine) {
	return func(e *Engine) {
		e.name = name
	}
}

// WithAuthor sets the author of the engine.
func WithAuthor(author string) func(*Engine) {
	return func(e *Engine) {
		e.author = author
	}
}

// WithLogger sets the logger used for the debug output.
func WithLogger(logger *log.Logger) func(*Engine) {
	return func(e *Engine) {
		e.logger = logger
	}
}

// WithNotation sets the chess notation.
func WithNotation(notation chess.Notation) func(*Engine) {
	return func(e *Engine) {
		e.notation = notation
	}
}

// WithSearch sets the search strategy.
func WithSearch(si search.Interface) func(*Engine) {
	return func(e *Engine) {
		e.options.search = si
	}
}

// WithEvaluation sets the evaluation strategy.
func WithEvaluation(ei evaluation.Interface) func(*Engine) {
	return func(e *Engine) {
		e.options.evaluation = ei
	}
}

// WithOracle sets the oracle strategy.
func WithOracle(oi oracle.Interface) func(*Engine) {
	return func(e *Engine) {
		e.options.oracle = oi
	}
}

// WithQuiescence sets the quiescence strategy.
func WithQuiescence(qi quiescence.Interface) func(*Engine) {
	return func(e *Engine) {
		e.options.quiescence = qi
	}
}

// WithTransposition sets the transposition strategy.
func WithTransposition(ti transposition.Interface) func(*Engine) {
	return func(e *Engine) {
		e.options.transposition = ti
	}
}

// WithOpening sets the opening strategy.
func WithOpening(oi opening.Interface) func(*Engine) {
	return func(e *Engine) {
		e.options.opening = oi
	}
}

// WithHash sets the size of the transposition hash table in MB.
func WithHash(hash int) func(*Engine) {
	return func(e *Engine) {
		e.options.hash = hash
	}
}

// Debug sets the debug option.
func (e *Engine) Debug(on bool) {
	e.debug = on
}

// Info returns the engine's info.
func (e *Engine) Info() (name, author string) {
	return e.name, e.author
}

// Init sets everything up.
func (e *Engine) Init() error {
	var err error
	e.once.Do(func() {
		err = e.options.transposition.Init(e.options.hash)
		e.initialized = true
	})
	return err
}

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
	if e.initialized {
		return errSetOption
	}

	for _, option := range availableOptions {
		if option.String() == name {
			fn, err := option.optionFunc(value)
			if err != nil {
				return err
			}
			fn(e)
			e.log("option", name, "set to", value)
			return nil
		}
	}

	return errOptionName
}

// SetPosition sets the position to the provided FEN.
func (e *Engine) SetPosition(fen string) error {
	fn, err := chess.FEN(fen)
	if err != nil {
		return err
	}
	fn(e.game)
	e.log("position set to", e.game.Position())
	return nil
}

// Move plays the moves on the current position.
func (e *Engine) Move(moves ...string) error {
	for _, move := range moves {
		m, err := e.notation.Decode(e.game.Position(), move)
		if err != nil {
			return err
		}

		if err := e.game.Move(m); err != nil {
			return err
		}
	}
	e.log("position set to", e.game.Position())
	return nil
}

// ResetPosition resets the position to the starting one.
func (e *Engine) ResetPosition() {
	e.game = chess.NewGame()
	e.log("position set to start")
}

// Search runs a search on the given input.
func (e *Engine) Search(ctx context.Context, input uci.Input) (<-chan uci.Output, error) {
	engineOutput := make(chan uci.Output)

	if !e.initialized {
		close(engineOutput)
		return engineOutput, errSearch
	}

	e.mu.Lock()
	start := time.Now()
	ctx, cancel := searchContext(ctx, input, e.stopSearch)

	searchMoves, err := searchMoves(e.notation, e.game.Position(), input.SearchMoves)
	if err != nil {
		e.log("could not parse search moves, defaulting to all possible moves", err)
	}
	searchOutput := search.Run(ctx, search.Input{
		Position:      e.game.Position(),
		SearchMoves:   searchMoves,
		Depth:         input.Depth,
		Search:        e.options.search,
		Evaluation:    e.options.evaluation,
		Oracle:        e.options.oracle,
		Quiescence:    e.options.quiescence,
		Transposition: e.options.transposition,
	})

	go func() {
		defer e.mu.Unlock()
		defer cancel()
		defer close(engineOutput)

		for output := range searchOutput {
			var pv []string
			for _, move := range output.PV {
				pv = append(pv, e.notation.Encode(e.game.Position(), move))
			}

			engineOutput <- uci.Output{
				Time:  time.Since(start),
				Depth: output.Depth,
				Nodes: output.Nodes,
				Score: output.Score,
				Mate:  output.Mate,
				PV:    pv,
			}
		}
	}()

	return engineOutput, nil
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

	e.options.transposition.Close()
}

// logger returns the logger to use depending on the debug setting
func (e *Engine) log(v ...any) {
	if e.debug {
		e.logger.Println(v...)
	}
}

// searchMoves decodes the moves to search from the engine notation
func searchMoves(notation chess.Notation, position *chess.Position, moves []string) ([]*chess.Move, error) {
	var next []*chess.Move
	for _, move := range moves {
		m, err := notation.Decode(position, move)
		if err != nil {
			return nil, err
		}
		next = append(next, m)
	}
	return next, nil
}

// searchContext creates a new context from the input
func searchContext(ctx context.Context, input uci.Input, stop <-chan struct{}) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

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
