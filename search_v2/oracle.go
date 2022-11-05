package search

import (
	"sort"

	"github.com/leonhfr/honeybadger/chess"
)

var promoOraclePoints = [13]int{0, 0, 5, 5, -5, -5, -5, -5, 10, 10, 0, 0, 0}

func orderMoves(moves []chess.Move) {
	sort.Slice(moves, func(i, j int) bool {
		return rank(moves[i]) > rank(moves[j])
	})
}

func rank(move chess.Move) (n int) {
	n += promoOraclePoints[move.Promo()]

	switch {
	case move.HasTag(chess.Capture):
		return 2
	case move.HasTag(chess.QueenSideCastle):
		return 3
	case move.HasTag(chess.KingSideCastle):
		return 4
	}

	return
}
