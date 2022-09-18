package engine

import (
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
			assert.Equal(t, tt.want.search, e.search)
		})
	}
}
