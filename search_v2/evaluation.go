// Package search implements the search algorithm.
package search

import "github.com/leonhfr/honeybadger/chess"

func isTerminal(pos *chess.Position, moves int) (int, bool) {
	switch {
	case moves == 0 && pos.InCheck():
		return -mate, true
	case moves == 0:
		return draw, true
	default:
		return 0, false
	}
}

func evaluate(pos *chess.Position) int {
	var mg, eg, phase int

	pos.PieceMap(func(p chess.Piece, sq chess.Square) {
		mgValue := pestoMGPieceTables[p][sq]
		egValue := pestoEGPieceTables[p][sq]

		if p.Color() == pos.Turn() {
			mg += mgValue
			eg += egValue
		} else {
			mg -= mgValue
			eg -= egValue
		}

		phase += pestoGamePhaseInc[p.Type()]
	})

	if phase > 24 {
		phase = 24 // in case of early promotion
	}

	return (phase*mg + (24-phase)*eg) / 24
}

// incMateDistance increase the distance to the mate by a count of one.
//
// In case of a positive score, it is decreased by 1.
// In case of a negative score, it is increased by 1.
func incMateDistance(score int) int {
	sign := 1
	if score < 0 {
		sign = -1
	}
	delta := mate - sign*score
	if delta <= maxDepth {
		return score - sign
	}
	return score
}
