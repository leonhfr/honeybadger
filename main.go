// Package main is the entry point of the engine.
package main

import (
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
	e := engine.New(fmt.Sprintf("%s %s", name, version), author)

	uci.Run(e, os.Stdin, os.Stdout)
}
