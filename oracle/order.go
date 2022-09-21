package oracle

import (
	"sort"

	"github.com/notnil/chess"
)

// Order implements move ordering based on available promotions and move tags.
// Priority is given first to queen and knight promotions. Them king side
// castling, queen side castling, checks, and captures.
type Order struct{}

// String implements the Interface interface.
func (Order) String() string {
	return "Order"
}

// Order implements the Interface interface.
func (Order) Order(moves []*chess.Move) {
	sort.Slice(moves, func(i, j int) bool {
		return rank(moves[i]) > rank(moves[j])
	})
}

func rank(move *chess.Move) (n int) {
	n += promoPoints[move.Promo()]

	for tag, points := range tagPoints {
		if move.HasTag(tag) {
			n += points
		}
	}

	return
}

// First, always consider promoting to a queen. Second, consider promoting to
// a knight for the edge cases when it is better. Third, consider the moves
// with no promotion at all. Last, consider promoting to the rest of the
// pieces. In no situation would promoting to one of these pieces be better
// than promoting to a queen or a knight.
var promoPoints = map[chess.PieceType]int{
	chess.Queen:       20,
	chess.Knight:      10,
	chess.NoPieceType: 0,
	chess.Rook:        -5,
	chess.Bishop:      -5,
	chess.Pawn:        -5,
	chess.King:        -5,
}

// After the promotions, we consider king side castling, then queen side
// castling, then checks, then captures.
var tagPoints = map[chess.MoveTag]int{
	chess.KingSideCastle:  5,
	chess.QueenSideCastle: 4,
	chess.Check:           3,
	chess.Capture:         2,
	chess.EnPassant:       1,
}
