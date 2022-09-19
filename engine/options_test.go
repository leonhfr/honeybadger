package engine

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/leonhfr/honeybadger/search"
	"github.com/leonhfr/honeybadger/uci"
)

func TestOptionStrategyString(t *testing.T) {
	assert.Equal(t, searchStrategy.name, searchStrategy.String())
}

func TestOptionStrategyUCI(t *testing.T) {
	assert.Equal(t, uci.Option{
		Type:    uci.OptionEnum,
		Name:    searchStrategy.name,
		Default: searchStrategy.def.String(),
		Vars: []string{
			search.Random{}.String(),
			search.Capture{}.String(),
			search.Negamax{}.String(),
			search.AlphaBeta{}.String(),
		},
	}, searchStrategy.uci())
}

// searchStrategy.defaultFunc tested in New

func TestOptionStrategyOptionFunc(t *testing.T) {
	type (
		args struct {
			search search.Interface
			value  string
		}
		want struct {
			search search.Interface
			err    error
		}
	)

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "option exists",
			args: args{search.Random{}, "Capture"},
			want: want{search.Capture{}, nil},
		},
		{
			name: "option does not exist",
			args: args{search.Random{}, ""},
			want: want{search.Random{}, errOptionValue},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn, err := searchStrategy.optionFunc(tt.args.value)
			assert.Equal(t, tt.want.err, err)

			e := New(WithSearch(tt.args.search))
			fn(e)
			assert.Equal(t, tt.want.search, e.options.search)
		})
	}
}

func TestOptionIntegerString(t *testing.T) {
	assert.Equal(t, hashOption.name, hashOption.String())
}

func TestOptionIntegerUCI(t *testing.T) {
	assert.Equal(t, uci.Option{
		Type:    uci.OptionInteger,
		Name:    hashOption.name,
		Default: fmt.Sprint(hashOption.def),
		Min:     fmt.Sprint(hashOption.min),
		Max:     fmt.Sprint(hashOption.max),
	}, hashOption.uci())
}

// optionInteger.defaultFunc tested in New

func TestOptionIntegerOptionFunc(t *testing.T) {
	type want struct {
		value int
		err   string
	}

	tests := []struct {
		name string
		args string
		want want
	}{
		{
			name: "value cannot be parsed as integer",
			args: "foobar",
			want: want{32, "strconv.ParseInt: parsing \"foobar\": invalid syntax"},
		},
		{
			name: "value is outside bounds",
			args: "0",
			want: want{32, errOutsideBound.Error()},
		},
		{
			name: "value is valid",
			args: "256",
			want: want{256, ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn, err := hashOption.optionFunc(tt.args)
			if err != nil {
				assert.Equal(t, tt.want.err, err.Error())
			}

			e := New()
			fn(e)
			assert.Equal(t, tt.want.value, e.options.hash)
		})
	}
}
