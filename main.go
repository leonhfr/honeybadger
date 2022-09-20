// Package main is the entry point of the engine.
package main

import (
	"context"
	"log"

	"github.com/leonhfr/honeybadger/cmd"
)

var (
	name    = "Honey Badger"
	version = "0.0.0"
	author  = "Leon Hollender"
)

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, cmd.NameKey, name)
	ctx = context.WithValue(ctx, cmd.VersionKey, version)
	ctx = context.WithValue(ctx, cmd.AuthorKey, author)

	if err := cmd.Execute(ctx); err != nil {
		log.Fatalf("honeybadger: %v", err)
	}
}
