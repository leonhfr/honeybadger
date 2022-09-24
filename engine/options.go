package engine

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/leonhfr/honeybadger/evaluation"
	"github.com/leonhfr/honeybadger/opening"
	"github.com/leonhfr/honeybadger/oracle"
	"github.com/leonhfr/honeybadger/quiescence"
	"github.com/leonhfr/honeybadger/search"
	"github.com/leonhfr/honeybadger/transposition"
	"github.com/leonhfr/honeybadger/uci"
)

var (
	errOptionValue  = errors.New("option value not found")
	errOutsideBound = errors.New("option value outside bounds")

	availableOptions = []option{
		searchStrategy,
		evaluationStrategy,
		oracleStrategy,
		quiescenceStrategy,
		transpositionStrategy,
		openingStrategy,
		hashOption,
	}

	searchStrategy = optionStrategy[search.Interface]{
		name: "SearchStrategy",
		def:  search.AlphaBeta{},
		vars: []search.Interface{
			search.Random{},
			search.Capture{},
			search.Negamax{},
			search.AlphaBeta{},
		},
		fn: WithSearch,
	}

	evaluationStrategy = optionStrategy[evaluation.Interface]{
		name: "EvaluationStrategy",
		def:  evaluation.Simplified{},
		vars: []evaluation.Interface{
			evaluation.Values{},
			evaluation.Simplified{},
		},
		fn: WithEvaluation,
	}

	oracleStrategy = optionStrategy[oracle.Interface]{
		name: "OracleStrategy",
		def:  oracle.Order{},
		vars: []oracle.Interface{
			oracle.None{},
			oracle.Order{},
		},
		fn: WithOracle,
	}

	quiescenceStrategy = optionStrategy[quiescence.Interface]{
		name: "QuiescenceStrategy",
		def:  quiescence.None{},
		vars: []quiescence.Interface{
			quiescence.None{},
			quiescence.AlphaBeta{},
		},
		fn: WithQuiescence,
	}

	transpositionStrategy = optionStrategy[transposition.Interface]{
		name: "TranspositionStrategy",
		def:  transposition.None{},
		vars: []transposition.Interface{
			transposition.None{},
			&transposition.Ristretto{},
		},
		fn: WithTransposition,
	}

	openingStrategy = optionStrategy[opening.Interface]{
		name: "OpeningStrategy",
		def:  opening.NewNone(),
		vars: []opening.Interface{
			opening.NewNone(),
			opening.NewBest(),
			opening.NewWeightedRandom(),
		},
		fn: WithOpening,
	}

	hashOption = optionInteger{
		name: "Hash",
		def:  32,
		min:  1,
		max:  1024,
		fn:   WithHash,
	}
)

// option is the interface implemented by each option type.
type option interface {
	fmt.Stringer
	uci() uci.Option
	defaultFunc() func(*Engine)
	optionFunc(value string) (func(*Engine), error)
}

// optionInteger represents an integer option.
type optionInteger struct {
	name          string
	def, min, max int
	fn            func(int) func(*Engine)
}

// String implements the option interface.
func (o optionInteger) String() string {
	return o.name
}

// uci implements the option interface.
func (o optionInteger) uci() uci.Option {
	return uci.Option{
		Type:    uci.OptionInteger,
		Name:    o.name,
		Default: fmt.Sprint(o.def),
		Min:     fmt.Sprint(o.min),
		Max:     fmt.Sprint(o.max),
	}
}

// defaultFunc implements the option interface.
func (o optionInteger) defaultFunc() func(*Engine) {
	return o.fn(o.def)
}

// optionFunc implements the option interface.
func (o optionInteger) optionFunc(value string) (func(*Engine), error) {
	v, err := strconv.ParseInt(value, 10, 0)
	if err != nil {
		return func(e *Engine) {}, err
	}

	if int(v) < o.min || int(v) > o.max {
		return func(e *Engine) {}, errOutsideBound
	}

	return o.fn(int(v)), nil
}

// optionStrategy represents a strategy option.
type optionStrategy[T fmt.Stringer] struct {
	name string
	def  T
	vars []T
	fn   func(T) func(*Engine)
}

// String implements the option interface.
func (o optionStrategy[T]) String() string {
	return o.name
}

// uci implements the option interface.
func (o optionStrategy[T]) uci() uci.Option {
	var vars []string
	for _, i := range o.vars {
		vars = append(vars, i.String())
	}
	return uci.Option{
		Type:    uci.OptionEnum,
		Name:    o.name,
		Default: o.def.String(),
		Vars:    vars,
	}
}

// defaultFunc implements the option interface.
func (o optionStrategy[T]) defaultFunc() func(*Engine) {
	return o.fn(o.def)
}

// optionFunc implements the option interface.
func (o optionStrategy[T]) optionFunc(value string) (func(*Engine), error) {
	for _, i := range o.vars {
		if value == i.String() {
			return o.fn(i), nil
		}
	}

	return func(e *Engine) {}, errOptionValue
}
