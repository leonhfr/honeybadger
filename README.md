# honeybadger

[![GoReportCard example](https://goreportcard.com/badge/github.com/leonhfr/honeybadger)](https://goreportcard.com/report/github.com/leonhfr/honeybadger)

Honey Badger is a UCI-compliant chess engine written in Go. Honey Badger is not a complete chess software and requires a UCI-compatible graphical user interface (GUI) to be used comfortably.

While being a toy project used for learning, it is working and actively maintained. Fair warning: it is not very strong.

Key features include:

- fully compliant UCI interface
- alpha-beta search with iterative deepening
- simple evaluation function combining piece values and positional advantage
- ability to use different search and evaluation strategies with options

Future (planned) features:

- quiescence search with null move pruning
- move ordering (oracle)
- better evaluation function with game phase knowledge
- transposition table for memoizing search results
- integrated opening book
- cli mode for quick analyses
- playable bot on Lichess

## Installation

Several installation methods are available:

- find the most recent [stable release](https://github.com/leonhfr/honeybadger/releases).
- using the `go` toolchain:

```sh
go install github.com/leonhfr/honeybadger@latest
```

- compile from source (requires `go@1.19` and `make`):

```sh
git clone git@github.com:leonhfr/honeybadger.git
cd honeybadger
make build
```

## Quick start

Honey Badger handles all of its communications via stdin and stdout using the UCI protocol. Therefore, a chess GUI that can communicate over UCI is needed. Refer to the documentation of your chosen GUI for information about how to use Honey Badger with it. We recommend:

- [leonhfr/cete](https://github.com/leonhfr/honeybadger), a CLI developed to pit UCI-compliant engines against each other. It runs games from command line options or configuration files, and can broadcast the game in a live web view.
- [cutechess/CuteChess](https://github.com/cutechess/cutechess)
- other options include SCID, Arena, Shredder, Fritz...

## Options

- **SearchStrategy**

  Search strategy to use. Available strategies are:

  - Random: plays random moves.
  - Capture: prioritizes capturing moves, and other plays random moves.
  - Negamax: implements the [negamax](https://en.wikipedia.org/wiki/Negamax) algorithm.
  - AlphaBeta (default): implements the negamax algorithm with [alpha-beta pruning](https://en.wikipedia.org/wiki/Alpha-beta_pruning).

- **EvaluationStrategy**

  Evaluation strategy to use. Available strategies are:

  - Values: difference between the piece values of each side.
  - Simplified (default): combination of piece values and positional advantage.

- **QuiescenceStrategy**

  Quiescence strategy to use. Available strategies are:

  - None (default): no quiescence search is performed.
  - AlphaBeta: negamax algorithm with alpha-beta pruning.
