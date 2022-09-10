package engine

import (
	"errors"
	"fmt"

	"github.com/leonhfr/honeybadger/search"
	"github.com/leonhfr/honeybadger/uci"
)

var (
	availableOptions = []option{
		searchStrategy,
	}

	searchStrategy = optionStrategy[search.Interface]{
		name: "SearchStrategy",
		def:  search.Capture{},
		vars: []search.Interface{
			search.Capture{},
			search.Random{},
		},
		fn: WithSearch,
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

	return func(e *Engine) {}, errors.New("option value not found")
}
