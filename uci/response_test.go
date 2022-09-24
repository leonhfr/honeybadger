package uci

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResponseString(t *testing.T) {
	tests := []struct {
		name string
		args response
		want string
	}{
		{name: "id", args: responseID{name: "NAME", author: "AUTHOR"}, want: "id name NAME\nid author AUTHOR"},
		{name: "uciok", args: responseUCIOK{}, want: "uciok"},
		{name: "readyok", args: responseReadyOK{}, want: "readyok"},
		{name: "bestmove", args: responseBestMove{"b1a3"}, want: "bestmove b1a3"},
		{
			name: "info score positive",
			args: Output{
				Depth: 8,
				Nodes: 1024,
				Score: 3000,
				PV:    []string{"b1a3", "b1c3"},
				Time:  time.Duration(5e9),
			},
			want: "info depth 8 nodes 1024 score cp 3000 pv b1a3 b1c3 time 5000",
		},
		{
			name: "info score negative",
			args: Output{
				Depth: 8,
				Nodes: 1024,
				Score: -3000,
				PV:    []string{"b1a3", "b1c3"},
				Time:  time.Duration(5e9),
			},
			want: "info depth 8 nodes 1024 score cp -3000 pv b1a3 b1c3 time 5000",
		},
		{
			name: "info mate positive",
			args: Output{
				Depth: 8,
				Score: 3000,
				Nodes: 1024,
				Mate:  5,
				PV:    []string{"b1a3", "b1c3"},
				Time:  time.Duration(5e9),
			},
			want: "info depth 8 nodes 1024 score mate 5 pv b1a3 b1c3 time 5000",
		},
		{
			name: "info mate negative",
			args: Output{
				Depth: 8,
				Nodes: 1024,
				Score: -3000,
				Mate:  -5,
				PV:    []string{"b1a3", "b1c3"},
				Time:  time.Duration(5e9),
			},
			want: "info depth 8 nodes 1024 score mate -5 pv b1a3 b1c3 time 5000",
		},
		{name: "comment", args: responseComment{comment: "COMMENT"}, want: "info string COMMENT"},
		{
			name: "option boolean",
			args: Option{
				Type:    OptionBoolean,
				Name:    "BOOLEAN OPTION",
				Default: "true",
			},
			want: "option name BOOLEAN OPTION type check default true",
		},
		{
			name: "option integer",
			args: Option{
				Type:    OptionInteger,
				Name:    "INTEGER OPTION",
				Default: "32",
				Min:     "2",
				Max:     "1024",
			},
			want: "option name INTEGER OPTION type spin default 32 min 2 max 1024",
		},
		{
			name: "option enum",
			args: Option{
				Type:    OptionEnum,
				Name:    "ENUM OPTION",
				Default: "Value1",
				Vars:    []string{"Value1", "Value2"},
			},
			want: "option name ENUM OPTION type combo default Value1 var Value1 var Value2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.String())
		})
	}
}
