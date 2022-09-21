package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"

	"github.com/leonhfr/honeybadger/engine"
	"github.com/leonhfr/honeybadger/uci"
)

const (
	depthFlag   = "depth"
	movesFlag   = "moves"
	timeFlag    = "time"
	verboseFlag = "verbose"
)

// searchCmd represents the search command.
var searchCmd = &cobra.Command{
	Use:   "search <fen>",
	Short: "Runs a single search on a FEN",
	Long:  `Search runs a single search on a FEN.`,
	Example: `  honeybadger search rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1 --time 10s
  honeybadger search 8/8/8/5K1k/8/8/8/5R2 w - - 0 1 --depth 3`,
	Args: cobra.ExactArgs(6),
	RunE: func(cmd *cobra.Command, args []string) error {
		fen := strings.Join(args, " ")

		e := engine.New()
		defer e.Quit()

		verbose := parseVerboseFlag(cmd)
		if verbose {
			e.Debug(true)
		}

		for _, option := range parseEngineOptionsFlags(cmd) {
			err := e.SetOption(option.name, option.value)
			if err != nil {
				return err
			}
		}

		if verbose {
			fmt.Println("initializing engine")
		}
		if err := e.Init(); err != nil {
			return err
		}

		if err := e.SetPosition(fen); err != nil {
			return err
		}

		input := parseInputFlags(cmd)
		if verbose {
			fmt.Println("running search...")
			fmt.Println(input)
		}
		oc, err := e.Search(cmd.Context(), input)
		if err != nil {
			return err
		}

		for output := range oc {
			fmt.Println(output)
		}

		return nil
	},
}

func init() {
	searchCmd.Flags().SortFlags = false

	// search options
	searchCmd.Flags().IntP(depthFlag, "d", 64, "depth at which to search")
	searchCmd.Flags().StringSliceP(movesFlag, "m", nil, "limit search to those moves in UCI notation")
	searchCmd.Flags().DurationP(timeFlag, "t", 0, "limit search time")
	searchCmd.Flags().Bool(verboseFlag, false, "verbose")

	// engine options
	for _, option := range engine.New().Options() {
		switch option.Type {
		case uci.OptionInteger:
			addIntegerOption(searchCmd, option)
		case uci.OptionEnum:
			addEnumOption(searchCmd, option)
		}
	}
}

func addIntegerOption(cmd *cobra.Command, option uci.Option) {
	value, _ := strconv.ParseInt(option.Default, 10, 0)
	cmd.Flags().Int(option.Name, int(value), fmt.Sprintf("from %s to %s", option.Min, option.Max))
}

func addEnumOption(cmd *cobra.Command, option uci.Option) {
	enum := newEnumFlag(option.Vars, option.Default)
	cmd.Flags().Var(enum, option.Name, fmt.Sprintf("one of %s", strings.Join(option.Vars, ", ")))
}

func parseInputFlags(cmd *cobra.Command) uci.Input {
	depth, _ := cmd.Flags().GetInt(depthFlag)
	moves, _ := cmd.Flags().GetStringSlice(movesFlag)
	time, _ := cmd.Flags().GetDuration(timeFlag)
	input := uci.Input{
		Depth:       depth,
		MoveTime:    time,
		SearchMoves: moves,
	}

	if time == 0 {
		input.Infinite = true
	}

	return input
}

func parseVerboseFlag(cmd *cobra.Command) bool {
	verbose, _ := cmd.Flags().GetBool(verboseFlag)
	return verbose
}

type engineOption struct{ name, value string }

func parseEngineOptionsFlags(cmd *cobra.Command) []engineOption {
	var options []engineOption
	for _, option := range engine.New().Options() {
		switch option.Type {
		case uci.OptionInteger:
			value, _ := cmd.Flags().GetInt(option.Name)
			options = append(options, engineOption{option.Name, fmt.Sprint(value)})
		case uci.OptionEnum:
			value, _ := cmd.Flags().GetString(option.Name)
			options = append(options, engineOption{option.Name, value})
		}
	}
	return options
}

type enumFlag struct {
	vars  []string
	value string
}

func newEnumFlag(vars []string, def string) *enumFlag {
	return &enumFlag{
		vars:  vars,
		value: def,
	}
}

func (ef *enumFlag) String() string {
	return ef.value
}

func (ef *enumFlag) Set(value string) error {
	if !slices.Contains(ef.vars, value) {
		return fmt.Errorf("%s should be one of %s", value, strings.Join(ef.vars, ", "))
	}
	ef.value = value
	return nil
}

func (ef *enumFlag) Type() string {
	return "string"
}
