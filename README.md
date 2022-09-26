# honeybadger

[![Lichess bot](https://img.shields.io/endpoint?style=flat&url=https%3A%2F%2Fpeppy-khapse-155e06.netlify.app%2F.netlify%2Ffunctions%2Fshield)](https://lichess.org/?user=honeybadger-bot#friend) [![GoReportCard](https://goreportcard.com/badge/github.com/leonhfr/honeybadger)](https://goreportcard.com/report/github.com/leonhfr/honeybadger)

Honey Badger is a UCI-compliant chess engine written in Go. Honey Badger is not a complete chess software and requires a UCI-compatible graphical user interface (GUI) to be used comfortably.

While it is a toy project used for learning, it's working and is actively maintained. Fair warning: it is not very strong.

Key features include:

- fully compliant UCI interface
- alpha-beta search with iterative deepening
- quiescence search
- oracle (move ordering)
- integrated opening book
- simple evaluation function combining piece values and positional advantage
- transposition table for memoizing search results
- ability to use different search and evaluation strategies with options
- cli mode for quick searches

Future (planned) features:

- null move pruning
- better evaluation function with game phase knowledge
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

## Quick start (UCI)

Honey Badger handles all of its communications via stdin and stdout using the [UCI protocol](https://backscattering.de/chess/uci/). Therefore, a chess GUI that can communicate over UCI is needed. Refer to the documentation of your chosen GUI for information about how to use Honey Badger with it. We recommend:

- [leonhfr/cete](https://github.com/leonhfr/honeybadger), a CLI developed to pit UCI-compliant engines against each other. It runs games from command line options or configuration files, and can broadcast the game in a live web view.
- [cutechess/CuteChess](https://github.com/cutechess/cutechess)
- other options include [SCID](http://scid.sourceforge.net/), [Arena](http://www.playwitharena.de/), [Shredder](https://www.shredderchess.com/)...

### Example

Executing the Honey Badger without any subcommands or arguments will run it in UCI mode:

```sh
honeybadger
# honeybadger now expects commands from stdin
```

Using [cete](https://github.com/leonhfr/honeybadger), you can quickly make UCI engines play games against each other using configuration files. For example:

```sh
# This will play a game between two Honey Badger, one playing randomly
# and the other using negamax with alpha-beta pruning
cete game ./test/data/cete/random-alphabeta.yaml


# This will play a game with the same options and will also broadcast
# the game in a web view
cete game -b ./test/data/cete/random-alphabeta.yaml

```

## Quick start (CLI)

Honey Badger has some limited features available from CLI subcommands.

```
Usage:
  honeybadger [flags]
  honeybadger [command]

Available Commands:
  help        Help about any command
  options     Lists the available options
  search      Runs a single search on a FEN

Flags:
  -h, --help      help for honeybadger
  -v, --version   version for honeybadger
```

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

- **OracleStrategy**

  Oracle strategy to use. Available strategies are:

  - None: no move ordering is performed.
  - Order (default): move ordering is performed based on promotions, castling, checks, and captures.

- **QuiescenceStrategy**

  Quiescence strategy to use. Available strategies are:

  - None (default): no quiescence search is performed.
  - AlphaBeta: negamax algorithm with alpha-beta pruning.

- **TranspositionStrategy**

  Transposition hash table strategy to use. Available strategies are:

  - None (default): no transposition hash table is used.
  - Ristretto: transposition hash table implemented using the [ristretto](https://github.com/dgraph-io/ristretto) library.

- **OpeningStrategy**

  Opening strategy to use. This defines how moves will be selected from the opening book. Available strategies are:

  - None (default): no opening strategy is used. The engine will not used any opening book and will only use the defined search strategy to determine which move to play.
  - Best: the best move from the opening book is played.
  - UniformRandom: moves have an equal probability of being chosen.
  - WeightedRandom: moves with a higher weight (quality) have a higher probability of being chosen.

- **Hash**

  Size of the transposition hash table in megabytes (MB).
  Defaults to 32 MB, can range from 1 to 1024 MB.
