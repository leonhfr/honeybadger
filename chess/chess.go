// Package chess provides types and functions to handle chess positions.
package chess

func pawnBitboard(sq Square, pos *Position) bitboard {
	bbSquare := sq.bitboard()
	var bbEnPassant bitboard
	if pos.enPassantSquare != NoSquare {
		bbEnPassant = pos.enPassantSquare.bitboard()
	}

	if pos.turn == White {
		captureR := (pos.board.bbBlack | bbEnPassant) & ((bbSquare & ^bbFileH & ^bbRank8) << 9)
		captureL := (pos.board.bbBlack | bbEnPassant) & ((bbSquare & ^bbFileA & ^bbRank8) << 7)
		upOne := pos.board.bbEmpty & ((bbSquare & ^bbRank8) << 8)
		upTwo := pos.board.bbEmpty & ((upOne & bbRank3) << 8)
		return captureR | captureL | upOne | upTwo
	}

	captureR := (pos.board.bbWhite | bbEnPassant) & ((bbSquare & ^bbFileH & ^bbRank1) >> 7)
	captureL := (pos.board.bbWhite | bbEnPassant) & ((bbSquare & ^bbFileA & ^bbRank1) >> 9)
	upOne := pos.board.bbEmpty & ((bbSquare & ^bbRank1) >> 8)
	upTwo := pos.board.bbEmpty & ((upOne & bbRank6) >> 8)
	return captureR | captureL | upOne | upTwo
}

func diagonalBitboard(sq Square, occupied bitboard) bitboard {
	square := sq.bitboard()
	return linearBitboard(square, occupied, bbDiagonals[sq]) |
		linearBitboard(square, occupied, bbAntiDiagonals[sq])
}

func hvBitboard(sq Square, occupied bitboard) bitboard {
	square := sq.bitboard()
	return linearBitboard(square, occupied, bbRanks[sq.Rank()/8]) |
		linearBitboard(square, occupied, bbFiles[sq.File()])
}

func linearBitboard(square, occupied, mask bitboard) bitboard {
	inMask := occupied & mask
	return ((inMask - 2*square) ^ (inMask.reverse() - 2*square.reverse()).reverse()) & mask
}
