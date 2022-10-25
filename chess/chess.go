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

		m := newMove(newPiece(check.color, King), NoPiece, check.s1, check.s2, pos.enPassant, NoPiece)
		moves = append(moves, m)
	}
	return moves
}

func standardMoves(pos *Position) []Move {
	bbAllowed := ^pos.bbWhite
	if pos.turn == Black {
		bbAllowed = ^pos.bbBlack
	}

	bbPinned := pos.getPinned(pos.turn)

	moves := []Move{}
	for p1 := newPiece(pos.turn, Pawn); p1 <= BlackKing; p1 += 2 {
		for bbS1 := pos.board.getBitboard(p1); bbS1 > 0; bbS1 = bbS1.resetLSB() {
			s1 := bbS1.scanForward()

			var bbS2 bitboard
			if bbPinned&s1.bitboard() > 0 {
				bbS2 = movePinnedBitboard(s1, pos, p1.Type()) & bbAllowed
			} else {
				bbS2 = moveBitboard(s1, pos, p1.Type()) & bbAllowed
			}

			for ; bbS2 > 0; bbS2 = bbS2.resetLSB() {
				s2 := bbS2.scanForward()

				p2 := pos.board.pieceByColor(s2, pos.turn.Other())
				if p1 == WhitePawn && s2.Rank() == Rank8 || p1 == BlackPawn && s2.Rank() == Rank1 {
					moves = append(moves,
						newMove(p1, p2, s1, s2, pos.enPassant, newPiece(pos.turn, Queen)),
						newMove(p1, p2, s1, s2, pos.enPassant, newPiece(pos.turn, Knight)),
						newMove(p1, p2, s1, s2, pos.enPassant, newPiece(pos.turn, Rook)),
						newMove(p1, p2, s1, s2, pos.enPassant, newPiece(pos.turn, Bishop)),
					)
				} else {
					moves = append(moves, newMove(p1, p2, s1, s2, pos.enPassant, NoPiece))
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
	ra := rookAttacksBitboard(sq, pos.bbOccupied)
	ba := bishopAttacksBitboard(sq, pos.bbOccupied)
	bb := bbKingMoves[sq] & pos.getBitboard(newPiece(c, King))
	bb |= (ra | ba) & pos.getBitboard(newPiece(c, Queen))
	bb |= ra & pos.getBitboard(newPiece(c, Rook))
	bb |= ba & pos.getBitboard(newPiece(c, Bishop))
	bb |= bbKnightMoves[sq] & pos.getBitboard(newPiece(c, Knight))

	if c == White {
		return (bb | bbBlackPawnCaptures[sq]&pos.bbWhitePawn) > 0
	}

	return (bb | bbWhitePawnCaptures[sq]&pos.bbBlackPawn) > 0
}

// isSquaresAttacked does not account for en passant attacks
func isSquaresAttacked(pos *Position, sqs ...Square) bool {
	c := pos.turn.Other()
	var bbR, bbB, bbK, bbN, bbP bitboard

	for _, sq := range sqs {
		bbR |= rookAttacksBitboard(sq, pos.bbOccupied)
		bbB |= bishopAttacksBitboard(sq, pos.bbOccupied)
		bbK |= bbKingMoves[sq]
		bbN |= bbKnightMoves[sq]

		if c == White {
			bbP |= bbBlackPawnCaptures[sq]
		} else {
			bbP |= bbWhitePawnCaptures[sq]
		}
	}

	bb := bbK & pos.getBitboard(newPiece(c, King))
	bb |= (bbR | bbB) & pos.getBitboard(newPiece(c, Queen))
	bb |= bbR & pos.getBitboard(newPiece(c, Rook))
	bb |= bbB & pos.getBitboard(newPiece(c, Bishop))
	bb |= bbN & pos.getBitboard(newPiece(c, Knight))
	bb |= bbP & pos.getBitboard(newPiece(c, Pawn))
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
		return ((bishopAttacksBitboard(sq, pos.bbOccupied) | rookAttacksBitboard(sq, pos.bbOccupied)) & bb).ones()
	case Rook:
		return (rookAttacksBitboard(sq, pos.bbOccupied) & bb).ones()
	case Bishop:
		return (bishopAttacksBitboard(sq, pos.bbOccupied) & bb).ones()
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
	if pos.enPassant != NoSquare {
		bbEnPassant = pos.enPassant.bitboard()
	}

	if pos.turn == White {
		captures := bbWhitePawnCaptures[sq]
		enPassantR := (bbSquare & (bbEnPassant << 8) & ^bbFileH) >> 1
		enPassantL := (bbSquare & (bbEnPassant << 8) & ^bbFileA) << 1
		return (pos.bbBlackPawn & (captures | enPassantR | enPassantL)).ones()
	}

	captures := bbBlackPawnCaptures[sq]
	enPassantR := (bbSquare & (bbEnPassant >> 8) & ^bbFileH) << 1
	enPassantL := (bbSquare & (bbEnPassant >> 8) & ^bbFileA) >> 1
	return (pos.bbWhitePawn & (captures | enPassantR | enPassantL)).ones()
}

func movePinnedBitboard(sq Square, pos *Position, pt PieceType) bitboard {
	king := pos.sqWhiteKing
	if pos.turn == Black {
		king = pos.sqBlackKing
	}

	switch pt {
	case Queen:
		return pinnedRookAttacksBitboard(sq, king, pos.bbOccupied) |
			pinnedBishopAttacksBitboard(sq, king, pos.bbOccupied)
	case Rook:
		return pinnedRookAttacksBitboard(sq, king, pos.bbOccupied)
	case Bishop:
		return pinnedBishopAttacksBitboard(sq, king, pos.bbOccupied)
	case Knight:
		return 0 // knights are always absolutely pinned
	case Pawn:
		pinner := pos.getPinner(pos.turn)
		bb := pinner & pawnCapturesBitboard(sq, pos)
		if bbFiles[sq]&pinner > 0 {
			bb |= pawnPushesBitboard(sq, pos)
		}
		return bb
	default:
		return 0
	}
}

func moveBitboard(sq Square, pos *Position, pt PieceType) bitboard {
	switch pt {
	case King:
		return bbKingMoves[sq]
	case Queen:
		return rookAttacksBitboard(sq, pos.bbOccupied) | bishopAttacksBitboard(sq, pos.bbOccupied)
	case Rook:
		return rookAttacksBitboard(sq, pos.bbOccupied)
	case Bishop:
		return bishopAttacksBitboard(sq, pos.bbOccupied)
	case Knight:
		return bbKnightMoves[sq]
	case Pawn:
		return pawnPushesBitboard(sq, pos) | pawnCapturesBitboard(sq, pos)
	default:
		return 0
	}
}

func pawnPushesBitboard(sq Square, pos *Position) bitboard {
	if pos.turn == White {
		upOne := ^pos.bbOccupied & bbWhitePawnPushes[sq]
		upTwo := ^pos.bbOccupied & ((upOne & bbRank3) << 8)
		return upOne | upTwo
	}

	upOne := ^pos.bbOccupied & bbBlackPawnPushes[sq]
	upTwo := ^pos.bbOccupied & ((upOne & bbRank6) >> 8)
	return upOne | upTwo
}

func pawnCapturesBitboard(sq Square, pos *Position) bitboard {
	var bbEnPassant bitboard
	if pos.enPassant != NoSquare {
		bbEnPassant = pos.enPassant.bitboard()
	}

	if pos.turn == White {
		return (pos.bbBlack | bbEnPassant) & bbWhitePawnCaptures[sq]
	}

	return (pos.bbWhite | bbEnPassant) & bbBlackPawnCaptures[sq]
}

func pinnedBitboard(sq Square, occupied, blockers, queen, rook, bishop bitboard) (bitboard, bitboard) {
	var pinned, pinner bitboard

	rPinner := xrayRookAttacksBitboard(sq, occupied, blockers) & (queen | rook)
	pinner |= rPinner

	for ; rPinner > 0; rPinner = rPinner.resetLSB() {
		s := rPinner.scanForward()
		pinned |= bbInBetween[sq][s] & blockers
	}

	bPinner := xrayBishopAttacksBitboard(sq, occupied, blockers) & (queen | bishop)
	pinner |= bPinner

	for ; bPinner > 0; bPinner = bPinner.resetLSB() {
		s := bPinner.scanForward()
		pinned |= bbInBetween[sq][s] & blockers
	}

	return pinned, pinner
}

func xrayBishopAttacksBitboard(sq Square, occupied, blockers bitboard) bitboard {
	attacks := bishopAttacksBitboard(sq, occupied)
	blockers &= attacks
	return attacks ^ bishopAttacksBitboard(sq, occupied^blockers)
}

func xrayRookAttacksBitboard(sq Square, occupied, blockers bitboard) bitboard {
	attacks := rookAttacksBitboard(sq, occupied)
	blockers &= attacks
	return attacks ^ rookAttacksBitboard(sq, occupied^blockers)
}

func pinnedBishopAttacksBitboard(sq, king Square, occupied bitboard) bitboard {
	if bb := bbDiagonals[sq]; bb&king.bitboard() > 0 {
		return linearBitboard(sq, occupied, bb)
	}

	if bb := bbAntiDiagonals[sq]; bb&king.bitboard() > 0 {
		return linearBitboard(sq, occupied, bb)
	}

	return 0
}

func pinnedRookAttacksBitboard(sq, king Square, occupied bitboard) bitboard {
	if bb := bbRanks[sq]; bb&king.bitboard() > 0 {
		return linearBitboard(sq, occupied, bb)
	}

	if bb := bbFiles[sq]; bb&king.bitboard() > 0 {
		return linearBitboard(sq, occupied, bb)
	}

	return 0
}

func bishopAttacksBitboard(sq Square, occupied bitboard) bitboard {
	return linearBitboard(sq, occupied, bbDiagonals[sq]) |
		linearBitboard(sq, occupied, bbAntiDiagonals[sq])
}

func rookAttacksBitboard(sq Square, occupied bitboard) bitboard {
	return linearBitboard(sq, occupied, bbRanks[sq]) |
		linearBitboard(sq, occupied, bbFiles[sq])
}

func linearBitboard(sq Square, occupied, mask bitboard) bitboard {
	inMask := occupied & mask
	return ((inMask - bbDoubleSquares[sq]) ^ (inMask.reverse() - bbReverseDoubleSquares[sq]).reverse()) & mask
}
