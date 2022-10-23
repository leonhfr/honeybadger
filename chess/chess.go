// Package chess provides types and functions to handle chess positions.
package chess

func pseudoMoves(pos *Position) []Move {
	return append(standardMoves(pos), castlingMoves(pos)...)
}

type castleCheck struct {
	color          Color
	side           Side
	s1, s2         Square
	travelBitboard bitboard
	checkSquares   []Square
}

var castleChecks = [2][2]castleCheck{
	{
		{White, KingSide, E1, G1, F1.bitboard() | G1.bitboard(), []Square{E1, F1, G1}},
		{White, QueenSide, E1, C1, B1.bitboard() | C1.bitboard() | D1.bitboard(), []Square{C1, D1, E1}},
	},
	{
		{Black, KingSide, E8, G8, F8.bitboard() | G8.bitboard(), []Square{E8, F8, G8}},
		{Black, QueenSide, E8, C8, B8.bitboard() | C8.bitboard() | D8.bitboard(), []Square{C8, D8, E8}},
	},
}

func castlingMoves(pos *Position) []Move {
	var moves []Move
	for _, check := range castleChecks[pos.turn] {
		if !pos.castlingRights.CanCastle(pos.turn, check.side) ||
			pos.bbOccupied&check.travelBitboard > 0 {
			continue
		}

		if isSquaresAttacked(pos, check.checkSquares...) {
			continue
		}

		m := newMove(newPiece(check.color, King), NoPiece, check.s1, check.s2, pos.enPassantSquare, NoPiece)
		moves = append(moves, m)
	}
	return moves
}

func standardMoves(pos *Position) []Move {
	bbAllowed := ^pos.bbWhite
	if pos.turn == Black {
		bbAllowed = ^pos.bbBlack
	}

	moves := []Move{}
	for _, p1 := range piecesByColor[pos.turn] {
		bbS1 := pos.board.getBitboard(p1)
		if bbS1 == 0 {
			continue
		}

		for s1 := A1; s1 <= H8; s1++ {
			if bbS1&s1.bitboard() == 0 {
				continue
			}

			bbS2 := moveBitboard(s1, pos, p1.Type()) & bbAllowed
			if bbS2 == 0 {
				continue
			}

			for s2 := A1; s2 <= H8; s2++ {
				if bbS2&s2.bitboard() == 0 {
					continue
				}

				p2 := pos.board.pieceByColor(s2, pos.turn.Other())
				if p1 == WhitePawn && s2.Rank() == Rank8 || p1 == BlackPawn && s2.Rank() == Rank1 {
					moves = append(moves,
						newMove(p1, p2, s1, s2, pos.enPassantSquare, newPiece(pos.turn, Queen)),
						newMove(p1, p2, s1, s2, pos.enPassantSquare, newPiece(pos.turn, Knight)),
						newMove(p1, p2, s1, s2, pos.enPassantSquare, newPiece(pos.turn, Rook)),
						newMove(p1, p2, s1, s2, pos.enPassantSquare, newPiece(pos.turn, Bishop)),
					)
				} else {
					moves = append(moves, newMove(p1, p2, s1, s2, pos.enPassantSquare, NoPiece))
				}
			}
		}
	}
	return moves
}

func isInCheck(pos *Position) bool {
	if pos.turn == White {
		return isAttacked(pos.sqWhiteKing, pos)
	}

	return isAttacked(pos.sqBlackKing, pos)
}

// isAttacked does not account for en passant attacks
func isAttacked(sq Square, pos *Position) bool {
	c := pos.turn.Other()
	hv := hvBitboard(sq, pos.bbOccupied)
	dia := diagonalBitboard(sq, pos.bbOccupied)
	r := bbKingMoves[sq] & pos.board.getBitboard(newPiece(c, King))
	r |= (hv | dia) & pos.board.getBitboard(newPiece(c, Queen))
	r |= hv & pos.board.getBitboard(newPiece(c, Rook))
	r |= dia & pos.board.getBitboard(newPiece(c, Bishop))
	r |= bbKnightMoves[sq] & pos.board.getBitboard(newPiece(c, Knight))

	if c == White {
		return (r | bbBlackPawnCaptures[sq]&pos.bbWhitePawn) > 0
	}

	return (r | bbWhitePawnCaptures[sq]&pos.bbBlackPawn) > 0
}

// isSquaresAttacked does not account for en passant attacks
func isSquaresAttacked(pos *Position, sqs ...Square) bool {
	c := pos.turn.Other()
	var bbHV, bbDia, bbK, bbN, bbP bitboard

	for _, sq := range sqs {
		hv := hvBitboard(sq, pos.bbOccupied)
		dia := diagonalBitboard(sq, pos.bbOccupied)
		bbHV |= hv
		bbDia |= dia
		bbK |= bbKingMoves[sq]
		bbN |= bbKnightMoves[sq]

		if c == White {
			bbP |= bbBlackPawnCaptures[sq]
		} else {
			bbP |= bbWhitePawnCaptures[sq]
		}
	}

	if c == White {
		bb := bbK & pos.bbWhiteKing
		bb |= (bbHV | bbDia) & pos.bbWhiteQueen
		bb |= bbHV & pos.bbWhiteRook
		bb |= bbDia & pos.bbWhiteBishop
		bb |= bbN & pos.bbWhiteKnight
		bb |= bbP & pos.bbWhitePawn
		return bb > 0
	}

	bb := bbK & pos.bbBlackKing
	bb |= (bbHV | bbDia) & pos.bbBlackQueen
	bb |= bbHV & pos.bbBlackRook
	bb |= bbDia & pos.bbBlackBishop
	bb |= bbN & pos.bbBlackKnight
	bb |= bbP & pos.bbBlackPawn
	return bb > 0
}

func isAttackedByCount(sq Square, pos *Position, by PieceType) int {
	switch bb := pos.getBitboard(newPiece(pos.turn.Other(), by)); by {
	case King:
		if (bb & bbKingMoves[sq]) != 0 {
			return 1
		}
		return 0
	case Queen:
		return ((diagonalBitboard(sq, pos.bbOccupied) | hvBitboard(sq, pos.bbOccupied)) & bb).ones()
	case Rook:
		return (hvBitboard(sq, pos.bbOccupied) & bb).ones()
	case Bishop:
		return (diagonalBitboard(sq, pos.bbOccupied) & bb).ones()
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
		return (pos.bbBlackPawn & (captures | enPassantR | enPassantL)).ones()
	}

	captures := bbBlackPawnCaptures[sq]
	enPassantR := (bbSquare & (bbEnPassant >> 8) & bbNotFileH) << 1
	enPassantL := (bbSquare & (bbEnPassant >> 8) & bbNotFileA) >> 1
	return (pos.bbWhitePawn & (captures | enPassantR | enPassantL)).ones()
}

func moveBitboard(sq Square, pos *Position, pt PieceType) bitboard {
	switch pt {
	case King:
		return bbKingMoves[sq]
	case Queen:
		return hvBitboard(sq, pos.bbOccupied) | diagonalBitboard(sq, pos.bbOccupied)
	case Rook:
		return hvBitboard(sq, pos.bbOccupied)
	case Bishop:
		return diagonalBitboard(sq, pos.bbOccupied)
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
		captures := (pos.bbBlack | bbEnPassant) & bbWhitePawnCaptures[sq]
		upOne := pos.bbEmpty & bbWhitePawnPushes[sq]
		upTwo := pos.bbEmpty & ((upOne & bbRank3) << 8)
		return captures | upOne | upTwo
	}

	captures := (pos.bbWhite | bbEnPassant) & bbBlackPawnCaptures[sq]
	upOne := pos.bbEmpty & bbBlackPawnPushes[sq]
	upTwo := pos.bbEmpty & ((upOne & bbRank6) >> 8)
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
