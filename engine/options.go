package engine

import (
	"errors"
	"fmt"

	"github.com/leonhfr/honeybadger/evaluation"
	"github.com/leonhfr/honeybadger/quiescence"
	"github.com/leonhfr/honeybadger/search"
	"github.com/leonhfr/honeybadger/uci"
)

var (
	errOptionName  = errors.New("option name not found")
	errOptionValue = errors.New("option value not found")

	availableOptions = []option{
		searchStrategy,
		evaluationStrategy,
		quiescenceStrategy,
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

	quiescenceStrategy = optionStrategy[quiescence.Interface]{
		name: "QuiescenceStrategy",
		def:  quiescence.None{},
		vars: []quiescence.Interface{
			quiescence.None{},
			quiescence.AlphaBeta{},
		},
		fn: WithQuiescence,
	}
)

// option is the interface implemented by each option type
type option interface {
	fmt.Stringer
	uci() uci.Option
	defaultFunc() func(*Engine)
	optionFunc(value string) (func(*Engine), error)
}

// optionStrategy represents a strategy option
type optionStrategy[T fmt.Stringer] struct {
	name string
	def  T
	vars []T
	fn   func(T) func(*Engine)
}

// String implements the option interface
func (o optionStrategy[T]) String() string {
	return o.name
}

// uci implements the option interface
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

// defaultFunc implements the option interface
func (o optionStrategy[T]) defaultFunc() func(*Engine) {
	return o.fn(o.def)
}

// optionFunc implements the option interface
func (o optionStrategy[T]) optionFunc(value string) (func(*Engine), error) {
	for _, i := range o.vars {
		if value == i.String() {
			return o.fn(i), nil
		}
	}

	return func(e *Engine) {}, errOptionValue
}
