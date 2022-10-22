// Package chess provides types and functions to handle chess positions.
package chess

func legalMoves(pos *Position) []Move {
	var moves []Move
	for _, m := range append(pseudoMoves(pos), castlingMoves(pos)...) {
		if pos.IsLegal(m) {
			moves = append(moves, m)
		}
	}
	return moves
}

type castleCheck struct {
	color          Color
	side           Side
	s1, s2         Square
	travelBitboard bitboard
	checkSquares   []Square
}

func (cc castleCheck) possible(pos *Position) bool {
	if pos.turn != cc.color ||
		!pos.castlingRights.CanCastle(pos.turn, cc.side) ||
		pos.board.bbOccupied&cc.travelBitboard > 0 {
		return false
	}

	for _, sq := range cc.checkSquares {
		if isAttacked(sq, pos) {
			return false
		}
	}

	return true
}

var castleChecks = [4]castleCheck{
	{White, KingSide, E1, G1, F1.bitboard() | G1.bitboard(), []Square{E1, F1, G1}},
	{White, QueenSide, E1, C1, B1.bitboard() | C1.bitboard() | D1.bitboard(), []Square{C1, D1, E1}},
	{Black, KingSide, E8, G8, F8.bitboard() | G8.bitboard(), []Square{E8, F8, G8}},
	{Black, QueenSide, E8, C8, B8.bitboard() | C8.bitboard() | D8.bitboard(), []Square{C8, D8, E8}},
}

func castlingMoves(pos *Position) []Move {
	moves := []Move{}
	for _, check := range castleChecks {
		if check.possible(pos) {
			m := newMove(newPiece(check.color, King), NoPiece, check.s1, check.s2, pos.enPassantSquare, NoPiece)
			moves = append(moves, m)
		}
	}
	return moves
}

var promoPieceTypes = [4]PieceType{Queen, Rook, Bishop, Knight}

func pseudoMoves(pos *Position) []Move {
	bbAllowed := ^pos.board.bbWhite
	if pos.turn == Black {
		bbAllowed = ^pos.board.bbBlack
	}

	moves := []Move{}
	for _, p1 := range piecesByColor[pos.turn] {
		s1bb := pos.board.getBitboard(p1)
		if s1bb == 0 {
			continue
		}

		for s1 := A1; s1 <= H8; s1++ {
			if s1bb&s1.bitboard() == 0 {
				continue
			}

			s2bb := moveBitboard(s1, pos, p1.Type()) & bbAllowed
			if s2bb == 0 {
				continue
			}

			for s2 := A1; s2 <= H8; s2++ {
				if s2bb&s2.bitboard() == 0 {
					continue
				}

				p2 := pos.board.piece(s2)
				if p1 == WhitePawn && s2.Rank() == Rank8 || p1 == BlackPawn && s2.Rank() == Rank1 {
					for _, pt := range promoPieceTypes {
						m := newMove(p1, p2, s1, s2, pos.enPassantSquare, newPiece(pos.turn, pt))
						moves = append(moves, m)
					}
				} else {
					m := newMove(p1, p2, s1, s2, pos.enPassantSquare, NoPiece)
					moves = append(moves, m)
				}
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
	for _, pt := range [6]PieceType{Queen, Rook, Bishop, Knight, Pawn, King} {
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
		return ((diagonalBitboard(sq, pos.board.bbOccupied) | hvBitboard(sq, pos.board.bbOccupied)) & bb).ones()
	case Rook:
		return (hvBitboard(sq, pos.board.bbOccupied) & bb).ones()
	case Bishop:
		return (diagonalBitboard(sq, pos.board.bbOccupied) & bb).ones()
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
		captures := bbWhitePawnCaptures[sq]
		enPassantR := (bbSquare & (bbEnPassant << 8) & bbNotFileH) >> 1
		enPassantL := (bbSquare & (bbEnPassant << 8) & bbNotFileA) << 1
		return (pos.board.bbBlackPawn & (captures | enPassantR | enPassantL)).ones()
	}

	captures := bbBlackPawnCaptures[sq]
	enPassantR := (bbSquare & (bbEnPassant >> 8) & bbNotFileH) << 1
	enPassantL := (bbSquare & (bbEnPassant >> 8) & bbNotFileA) >> 1
	return (pos.board.bbWhitePawn & (captures | enPassantR | enPassantL)).ones()
}

func moveBitboard(sq Square, pos *Position, pt PieceType) bitboard {
	switch pt {
	case King:
		return bbKingMoves[sq]
	case Queen:
		return hvBitboard(sq, pos.board.bbOccupied) | diagonalBitboard(sq, pos.board.bbOccupied)
	case Rook:
		return hvBitboard(sq, pos.board.bbOccupied)
	case Bishop:
		return diagonalBitboard(sq, pos.board.bbOccupied)
	case Knight:
		return bbKnightMoves[sq]
	case Pawn:
		return pawnBitboard(sq, pos)
	default:
		return 0
	}
}

func pawnBitboard(sq Square, pos *Position) bitboard {
	var bbEnPassant bitboard
	if pos.enPassantSquare != NoSquare {
		bbEnPassant = pos.enPassantSquare.bitboard()
	}

	if pos.turn == White {
		captures := (pos.board.bbBlack | bbEnPassant) & bbWhitePawnCaptures[sq]
		upOne := pos.board.bbEmpty & bbWhitePawnPushes[sq]
		upTwo := pos.board.bbEmpty & ((upOne & bbRank3) << 8)
		return captures | upOne | upTwo
	}

	captures := (pos.board.bbWhite | bbEnPassant) & bbBlackPawnCaptures[sq]
	upOne := pos.board.bbEmpty & bbBlackPawnPushes[sq]
	upTwo := pos.board.bbEmpty & ((upOne & bbRank6) >> 8)
	return captures | upOne | upTwo
}

func diagonalBitboard(sq Square, occupied bitboard) bitboard {
	return linearBitboard(sq, occupied, bbDiagonals[sq]) |
		linearBitboard(sq, occupied, bbAntiDiagonals[sq])
}

func hvBitboard(sq Square, occupied bitboard) bitboard {
	return linearBitboard(sq, occupied, bbRanks[sq]) |
		linearBitboard(sq, occupied, bbFiles[sq])
}

func linearBitboard(sq Square, occupied, mask bitboard) bitboard {
	inMask := occupied & mask
	return ((inMask - bbDoubleSquares[sq]) ^ (inMask.reverse() - bbReverseDoubleSquares[sq]).reverse()) & mask
}
