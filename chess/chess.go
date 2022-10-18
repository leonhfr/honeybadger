// Package chess provides types and functions to handle chess positions.
package chess

func diagonalMoves(sq Square, occupied bitboard) bitboard {
	square := sq.bitboard()
	return linearMoves(square, occupied, bbDiagonals[sq]) |
		linearMoves(square, occupied, bbAntiDiagonals[sq])
}

func hvMoves(sq Square, occupied bitboard) bitboard {
	square := sq.bitboard()
	return linearMoves(square, occupied, bbRanks[sq.Rank()/8]) |
		linearMoves(square, occupied, bbFiles[sq.File()])
}

func linearMoves(square, occupied, mask bitboard) bitboard {
	inMask := occupied & mask
	return ((inMask - 2*square) ^ (inMask.reverse() - 2*square.reverse()).reverse()) & mask
}
