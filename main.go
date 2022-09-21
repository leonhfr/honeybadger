// Package main is the entry point of the engine.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/leonhfr/honeybadger/cmd"
)

var (
	name    = "Honey Badger"
	version = "0.0.0"
	author  = "Leon Hollender"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ctx = context.WithValue(ctx, cmd.NameKey, name)
	ctx = context.WithValue(ctx, cmd.VersionKey, version)
	ctx = context.WithValue(ctx, cmd.AuthorKey, author)

	if err := cmd.Execute(ctx); err != nil {
		log.Fatalf("honeybadger: %v", err)
	}
}
