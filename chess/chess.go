// Package chess provides types and functions to handle chess positions.
package chess

var promoPieceTypes = []PieceType{Queen, Rook, Bishop, Knight}

func pseudoMoves(pos *Position) []*Move {
	bbAllowed := ^pos.board.bbWhite
	if pos.turn == Black {
		bbAllowed = ^pos.board.bbBlack
	}

	moves := []*Move{}
	for s1, p := range pos.board.squareMapByColor(pos.turn) {
		pseudoS2 := moveBitboard(s1, pos, p.Type()) & bbAllowed
		for s2 := range pseudoS2.mapping() {
			if p == WhitePawn && s2.Rank() == Rank8 || p == BlackPawn && s2.Rank() == Rank1 {
				for _, pt := range promoPieceTypes {
					moves = append(moves, newMove(pos, s1, s2, pt))
				}
			} else {
				moves = append(moves, newMove(pos, s1, s2, NoPieceType))
			}
		}
	}
	return moves
}

func isInCheck(pos *Position) bool {
	king := pos.board.sqWhiteKing
	if pos.turn == Black {
		king = pos.board.sqBlackKing
	}

	return isAttacked(king, pos)
}

func isAttacked(sq Square, pos *Position) bool {
	for _, pt := range []PieceType{Queen, Rook, Bishop, Knight, Pawn, King} {
		if isAttackedByCount(sq, pos, pt) > 0 {
			return true
		}
	}
	return false
}

func isAttackedByCount(sq Square, pos *Position, by PieceType) int {
	switch bb := pos.board.getBitboard(newPiece(pos.turn.Other(), by)); by {
	case King:
		if (bb & bbKingMoves[sq]) != 0 {
			return 1
		}
		return 0
	case Queen:
		return ((diagonalBitboard(sq, ^pos.board.bbEmpty) | hvBitboard(sq, ^pos.board.bbEmpty)) & bb).ones()
	case Rook:
		return (hvBitboard(sq, ^pos.board.bbEmpty) & bb).ones()
	case Bishop:
		return (diagonalBitboard(sq, ^pos.board.bbEmpty) & bb).ones()
	case Knight:
		return (bbKnightMoves[sq] & bb).ones()
	case Pawn:
		return isAttackedByPawnCount(sq, pos)
	default:
		return 0
	}
}

func isAttackedByPawnCount(sq Square, pos *Position) int {
	bbSquare := sq.bitboard()
	var bbEnPassant bitboard
	if pos.enPassantSquare != NoSquare {
		bbEnPassant = pos.enPassantSquare.bitboard()
	}

	if pos.turn == White {
		captureR := (bbSquare & ^bbFileH & ^bbRank8) << 9
		captureL := (bbSquare & ^bbFileA & ^bbRank8) << 7
		enPassantR := (bbSquare & (bbEnPassant << 8) & ^bbFileH) >> 1
		enPassantL := (bbSquare & (bbEnPassant << 8) & ^bbFileA) << 1
		return (pos.board.bbBlackPawn & (captureR | captureL | enPassantR | enPassantL)).ones()
	}

	captureR := (bbSquare & ^bbFileH & ^bbRank1) >> 7
	captureL := (bbSquare & ^bbFileA & ^bbRank1) >> 9
	enPassantR := (bbSquare & (bbEnPassant >> 8) & ^bbFileH) << 1
	enPassantL := (bbSquare & (bbEnPassant >> 8) & ^bbFileA) >> 1
	return (pos.board.bbWhitePawn & (captureR | captureL | enPassantR | enPassantL)).ones()
}

func moveBitboard(sq Square, pos *Position, pt PieceType) bitboard {
	switch pt {
	case King:
		return bbKingMoves[sq]
	case Queen:
		return hvBitboard(sq, ^pos.board.bbEmpty) | diagonalBitboard(sq, ^pos.board.bbEmpty)
	case Rook:
		return hvBitboard(sq, ^pos.board.bbEmpty)
	case Bishop:
		return diagonalBitboard(sq, ^pos.board.bbEmpty)
	case Knight:
		return bbKnightMoves[sq]
	case Pawn:
		return pawnBitboard(sq, pos)
	default:
		return 0
	}
}

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
