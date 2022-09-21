// Package cmd implements the different commands.
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/leonhfr/honeybadger/engine"
	"github.com/leonhfr/honeybadger/uci"
)

// ContextKey represents a value's key in the context.
type ContextKey int

const (
	NameKey    ContextKey = iota // NameKey is the key to the name value.
	VersionKey                   // VersionKey is the key to the version value.
	AuthorKey                    // AuthorKey is the key to the author value.
)

// rootCmd represents the base command when called without any subcommands.
// It is the entry point of the UCI interface.
var rootCmd = &cobra.Command{
	Use:   "honeybadger",
	Short: "Honey Badger is a UCI-compliant chess engine written in Go",
	Long: `Honey Badger is a UCI-compliant chess engine written in Go.

Honey Badger is not a complete chess software and requires a UCI-compatible
graphical user interface (GUI) to be used comfortably. The root command
is the entrypoint of the UCI interface and will accept commands from stdin
and will output responses in stdout.

Fair warning: it is not very strong.`,
	Args:              cobra.MatchAll(cobra.NoArgs),
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		name, version, author := name(ctx), version(ctx), author(ctx)

		e := engine.New(
			engine.WithName(fmt.Sprintf("%s v%s", name, version)),
			engine.WithAuthor(author),
			engine.WithLogger(uci.Logger(os.Stdout)),
		)

		uci.Run(ctx, e, os.Stdin, os.Stdout)

		return nil
	},
}

// Execute is the entry point of the root command and subcommands.
func Execute(ctx context.Context) error {
	rootCmd.Version = version(ctx)

	return rootCmd.ExecuteContext(ctx)
}

func init() {
	rootCmd.AddCommand(optionsCmd)
}

// name returns the name value from the context.
func name(ctx context.Context) string {
	name, ok := ctx.Value(NameKey).(string)
	if !ok {
		return ""
	}
	return name
}

// version returns the version value from the context.
func version(ctx context.Context) string {
	version, ok := ctx.Value(VersionKey).(string)
	if !ok {
		return "0.0.0"
	}
	return version
}

// author returns the author value from the context.
func author(ctx context.Context) string {
	author, ok := ctx.Value(AuthorKey).(string)
	if !ok {
		return ""
	}
	return author
}
