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
		pieceValue := simplifiedPieceValues[piece.Type()]
		var values [64]int

		for i := 0; i < 64; i++ {
			positionValueIndex := i // black
			if piece.Color() == chess.White {
				positionValueIndex = i + 56 - 2*8*(i/8)
			}

			values[i] = pieceValue + simplifiedHumanPieceTables[piece.Type()][positionValueIndex]
		}

		simplifiedPieceTables[piece] = values
	}
}

var (
	simplifiedPieceTables = make(map[chess.Piece][64]int)

	simplifiedPieceValues = map[chess.PieceType]int{
		chess.King:   20000,
		chess.Queen:  900,
		chess.Rook:   500,
		chess.Bishop: 330,
		chess.Knight: 320,
		chess.Pawn:   100,
	}

	simplifiedHumanPieceTables = map[chess.PieceType][64]int{
		chess.King:   simplifiedKingTable,
		chess.Queen:  simplifiedQueenTable,
		chess.Rook:   simplifiedRookTable,
		chess.Bishop: simplifiedBishopTable,
		chess.Knight: simplifiedKnightTable,
		chess.Pawn:   simplifiedPawnTable,
	}

	simplifiedPawnTable = [64]int{
		0, 0, 0, 0, 0, 0, 0, 0,
		50, 50, 50, 50, 50, 50, 50, 50,
		10, 10, 20, 30, 30, 20, 10, 10,
		5, 5, 10, 25, 25, 10, 5, 5,
		0, 0, 0, 20, 20, 0, 0, 0,
		5, -5, -10, 0, 0, -10, -5, 5,
		5, 10, 10, -20, -20, 10, 10, 5,
		0, 0, 0, 0, 0, 0, 0, 0,
	}

	simplifiedKnightTable = [64]int{
		-50, -40, -30, -30, -30, -30, -40, -50,
		-40, -20, 0, 0, 0, 0, -20, -40,
		-30, 0, 10, 15, 15, 10, 0, -30,
		-30, 5, 15, 20, 20, 15, 5, -30,
		-30, 0, 15, 20, 20, 15, 0, -30,
		-30, 5, 10, 15, 15, 10, 5, -30,
		-40, -20, 0, 5, 5, 0, -20, -40,
		-50, -40, -30, -30, -30, -30, -40, -50,
	}

	simplifiedBishopTable = [64]int{
		-20, -10, -10, -10, -10, -10, -10, -20,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-10, 0, 5, 10, 10, 5, 0, -10,
		-10, 5, 5, 10, 10, 5, 5, -10,
		-10, 0, 10, 10, 10, 10, 0, -10,
		-10, 10, 10, 10, 10, 10, 10, -10,
		-10, 5, 0, 0, 0, 0, 5, -10,
		-20, -10, -10, -10, -10, -10, -10, -20,
	}

	simplifiedRookTable = [64]int{
		0, 0, 0, 0, 0, 0, 0, 0,
		5, 10, 10, 10, 10, 10, 10, 5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		0, 0, 0, 5, 5, 0, 0, 0,
	}

	simplifiedQueenTable = [64]int{
		-20, -10, -10, -5, -5, -10, -10, -20,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-10, 0, 5, 5, 5, 5, 0, -10,
		-5, 0, 5, 5, 5, 5, 0, -5,
		0, 0, 5, 5, 5, 5, 0, -5,
		-10, 5, 5, 5, 5, 5, 0, -10,
		-10, 0, 5, 0, 0, 0, 0, -10,
		-20, -10, -10, -5, -5, -10, -10, -20,
	}

	simplifiedKingTable = [64]int{
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-20, -30, -30, -40, -40, -30, -30, -20,
		-10, -20, -20, -20, -20, -20, -20, -10,
		20, 20, 0, 0, 0, 0, 20, 20,
		20, 30, 10, 0, 0, 10, 30, 20,
	}
)
