// Package main is the entry point of the engine.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/leonhfr/honeybadger/engine"
	"github.com/leonhfr/honeybadger/uci"
)

var (
	name    = "Honey Badger"
	version = "0.0.0"
	author  = "Leon Hollender"
)

func main() {
	e := engine.New(
		engine.WithName(fmt.Sprintf("%s v%s", name, version)),
		engine.WithAuthor(author),
		engine.WithLogger(uci.Logger(os.Stdout)),
	)

	uci.Run(context.Background(), e, os.Stdin, os.Stdout)
}
