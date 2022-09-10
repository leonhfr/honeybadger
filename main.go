// Package main is the entry point of the engine.
package main

import (
	"fmt"
	"os"

	"github.com/leonhfr/honeybadger/engine"
	"github.com/leonhfr/honeybadger/uci"
)

var (
	engineName    = "Honey Badger"
	engineVersion = "0.0.0"
	engineAuthor  = "Leon Hollender"
)

func main() {
	e := engine.New(engineName, fmt.Sprintf("%s %s", engineAuthor, engineVersion))

	uci.Run(e, os.Stdin, os.Stdout)
}
