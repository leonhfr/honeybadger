package evaluation

import "github.com/notnil/chess"

var pieceValues = map[chess.PieceType]int{
	chess.King:        20000,
	chess.Queen:       800,
	chess.Rook:        500,
	chess.Bishop:      300,
	chess.Knight:      300,
	chess.Pawn:        100,
	chess.NoPieceType: 0,
}

// Values simply subtracts all piece values from each side.
type Values struct{}

// String implements the Interface interface.
func (Values) String() string {
	return "Values"
}

// Evaluate implements the Interface interface.
func (Values) Evaluate(p *chess.Position) int {
	var value int
	for _, piece := range p.Board().SquareMap() {
		if piece.Color() == p.Turn() {
			value += pieceValues[piece.Type()]
		} else {
			value -= pieceValues[piece.Type()]
		}
	}
	return value
}
