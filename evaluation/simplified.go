package evaluation

import (
	"github.com/notnil/chess"
)

// Simplified implements the evaluation function from Tomasz Michniewski.
// Source: https://www.chessprogramming.org/Simplified_Evaluation_Function
type Simplified struct{}

// String implements the Interface interface.
func (Simplified) String() string {
	return "Simplified"
}

// Evaluate implements the Interface interface.
func (Simplified) Evaluate(p *chess.Position) int {
	var value int
	for square, piece := range p.Board().SquareMap() {
		pieceValue := simplifiedPieceTables[piece][int(square)]

		if piece.Color() == p.Turn() {
			value += pieceValue
		} else {
			value -= pieceValue
		}
	}
	return value
}

func init() {
	for _, piece := range []chess.Piece{
		chess.WhiteKing, chess.WhiteQueen, chess.WhiteRook,
		chess.WhiteBishop, chess.WhiteKnight, chess.WhitePawn,
		chess.BlackKing, chess.BlackQueen, chess.BlackRook,
		chess.BlackBishop, chess.BlackKnight, chess.BlackPawn,
	} {
		value := simplifiedPieceValues[piece.Type()]
		human := simplifiedHumanPieceTables[piece.Type()]

		if piece.Color() == chess.White {
			simplifiedPieceTables[piece] = mapSquareTableToWhite(human, value)
		} else {
			simplifiedPieceTables[piece] = mapSquareTableToBlack(human, value)
		}
	}
}

func mapSquareTableToWhite(human [8][8]int, value int) [64]int {
	var table [64]int
	for i := 0; i < 64; i++ {
		table[i] = value + human[7-i/8][i%8]
	}
	return table
}

func mapSquareTableToBlack(human [8][8]int, value int) [64]int {
	var table [64]int
	for i := 0; i < 64; i++ {
		table[i] = value + human[i/8][i%8]
	}
	return table
}

var (
	simplifiedPieceTables = make(map[chess.Piece][64]int)

	simplifiedPieceValues = map[chess.PieceType]int{
		chess.King:   0,
		chess.Queen:  900,
		chess.Rook:   500,
		chess.Bishop: 330,
		chess.Knight: 320,
		chess.Pawn:   100,
	}

	simplifiedHumanPieceTables = map[chess.PieceType][8][8]int{
		chess.King:   simplifiedHumanKingTable,
		chess.Queen:  simplifiedHumanQueenTable,
		chess.Rook:   simplifiedHumanRookTable,
		chess.Bishop: simplifiedHumanBishopTable,
		chess.Knight: simplifiedHumanKnightTable,
		chess.Pawn:   simplifiedHumanPawnTable,
	}

	simplifiedHumanPawnTable = [8][8]int{
		{0, 0, 0, 0, 0, 0, 0, 0},
		{50, 50, 50, 50, 50, 50, 50, 50},
		{10, 10, 20, 30, 30, 20, 10, 10},
		{5, 5, 10, 25, 25, 10, 5, 5},
		{0, 0, 0, 20, 20, 0, 0, 0},
		{5, -5, -10, 0, 0, -10, -5, 5},
		{5, 10, 10, -20, -20, 10, 10, 5},
		{0, 0, 0, 0, 0, 0, 0, 0},
	}

	simplifiedHumanKnightTable = [8][8]int{
		{-50, -40, -30, -30, -30, -30, -40, -50},
		{-40, -20, 0, 0, 0, 0, -20, -40},
		{-30, 0, 10, 15, 15, 10, 0, -30},
		{-30, 5, 15, 20, 20, 15, 5, -30},
		{-30, 0, 15, 20, 20, 15, 0, -30},
		{-30, 5, 10, 15, 15, 10, 5, -30},
		{-40, -20, 0, 5, 5, 0, -20, -40},
		{-50, -40, -30, -30, -30, -30, -40, -50},
	}

	simplifiedHumanBishopTable = [8][8]int{
		{-20, -10, -10, -10, -10, -10, -10, -20},
		{-10, 0, 0, 0, 0, 0, 0, -10},
		{-10, 0, 5, 10, 10, 5, 0, -10},
		{-10, 5, 5, 10, 10, 5, 5, -10},
		{-10, 0, 10, 10, 10, 10, 0, -10},
		{-10, 10, 10, 10, 10, 10, 10, -10},
		{-10, 5, 0, 0, 0, 0, 5, -10},
		{-20, -10, -10, -10, -10, -10, -10, -20},
	}

	simplifiedHumanRookTable = [8][8]int{
		{0, 0, 0, 0, 0, 0, 0, 0},
		{5, 10, 10, 10, 10, 10, 10, 5},
		{-5, 0, 0, 0, 0, 0, 0, -5},
		{-5, 0, 0, 0, 0, 0, 0, -5},
		{-5, 0, 0, 0, 0, 0, 0, -5},
		{-5, 0, 0, 0, 0, 0, 0, -5},
		{-5, 0, 0, 0, 0, 0, 0, -5},
		{0, 0, 0, 5, 5, 0, 0, 0},
	}

	simplifiedHumanQueenTable = [8][8]int{
		{-20, -10, -10, -5, -5, -10, -10, -20},
		{-10, 0, 0, 0, 0, 0, 0, -10},
		{-10, 0, 5, 5, 5, 5, 0, -10},
		{-5, 0, 5, 5, 5, 5, 0, -5},
		{0, 0, 5, 5, 5, 5, 0, -5},
		{-10, 5, 5, 5, 5, 5, 0, -10},
		{-10, 0, 5, 0, 0, 0, 0, -10},
		{-20, -10, -10, -5, -5, -10, -10, -20},
	}

	simplifiedHumanKingTable = [8][8]int{
		{-30, -40, -40, -50, -50, -40, -40, -30},
		{-30, -40, -40, -50, -50, -40, -40, -30},
		{-30, -40, -40, -50, -50, -40, -40, -30},
		{-30, -40, -40, -50, -50, -40, -40, -30},
		{-20, -30, -30, -40, -40, -30, -30, -20},
		{-10, -20, -20, -20, -20, -20, -20, -10},
		{20, 20, 0, 0, 0, 0, 20, 20},
		{20, 30, 10, 0, 0, 10, 30, 20},
	}
)
