package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	startingSquareMap = SquareMap{
		A8: BlackRook, B8: BlackKnight, C8: BlackBishop, D8: BlackQueen,
		E8: BlackKing, F8: BlackBishop, G8: BlackKnight, H8: BlackRook,
		A7: BlackPawn, B7: BlackPawn, C7: BlackPawn, D7: BlackPawn,
		E7: BlackPawn, F7: BlackPawn, G7: BlackPawn, H7: BlackPawn,
		A2: WhitePawn, B2: WhitePawn, C2: WhitePawn, D2: WhitePawn,
		E2: WhitePawn, F2: WhitePawn, G2: WhitePawn, H2: WhitePawn,
		A1: WhiteRook, B1: WhiteKnight, C1: WhiteBishop, D1: WhiteQueen,
		E1: WhiteKing, F1: WhiteBishop, G1: WhiteKnight, H1: WhiteRook,
	}

	startingBoard = board{
		bbWhiteKing:   16,
		bbWhiteQueen:  8,
		bbWhiteRook:   129,
		bbWhiteBishop: 36,
		bbWhiteKnight: 66,
		bbWhitePawn:   65280,
		bbBlackKing:   1152921504606846976,
		bbBlackQueen:  576460752303423488,
		bbBlackRook:   9295429630892703744,
		bbBlackBishop: 2594073385365405696,
		bbBlackKnight: 4755801206503243776,
		bbBlackPawn:   71776119061217280,
	}
)

func TestNewBoard(t *testing.T) {
	assert.Equal(t, startingBoard, *newBoard(startingSquareMap))
}

func TestBoard_SquareMap(t *testing.T) {
	assert.Equal(t, startingSquareMap, startingBoard.squareMap())
}
