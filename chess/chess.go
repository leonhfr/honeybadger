// Package chess provides types and functions to handle chess positions.
package chess

func pseudoMoves(pos *Position) []Move {
	switch checks := pos.getCheck(pos.turn.Other()); checks.ones() {
	case 0: // no check
		return append(standardMoves(pos), castlingMoves(pos)...)
	case 1: // single check
		return append(checkAttackAndInterposingMoves(pos), checkFlightMoves(pos)...)
	default: // double check and more
		return checkFlightMoves(pos)
	}
}

// assumes there is only one checking piece
func checkAttackAndInterposingMoves(pos *Position) []Move {
	c := pos.turn
	bbChecking := pos.getCheck(c.Other())
	sqKing, sqChecking := pos.getKingSquare(c), bbChecking.scanForward()
	bbBetween := ^pos.getColor(c) & bbInBetween[sqKing][sqChecking]
	bbPinned := pos.getPinned(c)

	var moves []Move
	for p1 := Pawn.color(c); p1 <= WhiteQueen; p1 += 2 {
		for bbS1 := pos.getBitboard(p1) & ^bbPinned; bbS1 > 0; bbS1 = bbS1.resetLSB() {
			s1 := bbS1.scanForward()
			bbS2 := moveBitboard(s1, pos, p1.Type())

			bbTarget := bbChecking
			if p1.Type() == Pawn && pos.pieceAt(sqChecking).Type() == Pawn {
				bbTarget |= pos.enPassant.bitboard()
			}

			// attacks to the attacking piece (not pinned)
			for bbAttack := bbS2 & bbTarget; bbAttack > 0; bbAttack = bbAttack.resetLSB() {
				s2 := bbAttack.scanForward()
				p2 := pos.pieceAt(s2)

				if p1 == WhitePawn && s2.Rank() == Rank8 || p1 == BlackPawn && s2.Rank() == Rank1 {
					moves = append(moves,
						newMove(p1, p2, s1, s2, pos.enPassant, Queen.color(c)),
						newMove(p1, p2, s1, s2, pos.enPassant, Rook.color(c)),
						newMove(p1, p2, s1, s2, pos.enPassant, Bishop.color(c)),
						newMove(p1, p2, s1, s2, pos.enPassant, Knight.color(c)),
					)
				} else {
					moves = append(moves, newMove(p1, p2, s1, s2, pos.enPassant, NoPiece))
				}
			}

			// interposing moves in case of distance sliding checks (not pinned)
			for bbInterpose := bbS2 & bbBetween; bbInterpose > 0; bbInterpose = bbInterpose.resetLSB() {
				s2 := bbInterpose.scanForward()

				// there is always a clear path so p2 is always NoPiece
				moves = append(moves, newMove(p1, NoPiece, s1, s2, NoSquare, NoPiece))
			}
		}
	}
	return moves
}

// king moves to non attacked squares
func checkFlightMoves(pos *Position) []Move {
	c, op := pos.turn, pos.turn.Other()

	bbFlight := bbKingMoves[pos.getKingSquare(c)]
	bbFlight &= ^pos.getColor(c)                    // possible moves
	bbFlight &= ^bbKingMoves[pos.getKingSquare(op)] // enemy king attacks

	// pawn attacks
	if c == White {
		bbFlight &= ^((pos.bbBlack & pos.bbPawn & ^bbFileA) >> 9)
		bbFlight &= ^((pos.bbBlack & pos.bbPawn & ^bbFileH) >> 7)
	} else {
		bbFlight &= ^((pos.bbWhite & pos.bbPawn &^ bbFileA) << 9)
		bbFlight &= ^((pos.bbWhite & pos.bbPawn &^ bbFileH) << 7)
	}

	for bb := pos.getBitboard(Knight.color(op)); bb > 0; bb = bb.resetLSB() {
		sq := bb.scanForward()
		bbFlight &= ^bbKnightMoves[sq]
	}

	bbOccupiedNoKing := pos.bbOccupied & ^pos.getBitboard(King.color(c))

	for bb := pos.getBitboard(Queen.color(op)); bb > 0; bb = bb.resetLSB() {
		sq := bb.scanForward()
		bbFlight &= ^rookAttacksBitboard(sq, bbOccupiedNoKing)
		bbFlight &= ^bishopAttacksBitboard(sq, bbOccupiedNoKing)
	}

	for bb := pos.getBitboard(Rook.color(op)); bb > 0; bb = bb.resetLSB() {
		sq := bb.scanForward()
		bbFlight &= ^rookAttacksBitboard(sq, bbOccupiedNoKing)
	}

	for bb := pos.getBitboard(Bishop.color(op)); bb > 0; bb = bb.resetLSB() {
		sq := bb.scanForward()
		bbFlight &= ^bishopAttacksBitboard(sq, bbOccupiedNoKing)
	}

	var moves []Move
	for s1 := pos.getKingSquare(c); bbFlight > 0; bbFlight = bbFlight.resetLSB() {
		s2 := bbFlight.scanForward()
		p1, p2 := King.color(c), pos.pieceAt(s2)
		moves = append(moves, newMove(p1, p2, s1, s2, NoSquare, NoPiece))
	}
	return moves
}

const (
	bbKingCastleTravel  bitboard = 1<<F1 | 1<<G1
	bbQueenCastleTravel bitboard = 1<<B1 | 1<<C1 | 1<<D1
)

func castlingMoves(pos *Position) []Move {
	var moves []Move

	c := pos.turn
	bbKingTravel, bbQueenTravel := bbKingCastleTravel, bbQueenCastleTravel
	sqKing, sqKingSide, sqQueenSide := E1, G1, C1
	if c == Black {
		bbKingTravel <<= 56
		bbQueenTravel <<= 56
		sqKing += 56
		sqKingSide += 56
		sqQueenSide += 56
	}

	if pos.castlingRights.CanCastle(c, KingSide) &&
		pos.bbOccupied&bbKingTravel == 0 {
		moves = append(moves, newMove(
			King.color(c),
			NoPiece,
			sqKing,
			sqKingSide,
			NoSquare,
			NoPiece,
		))
	}

	if pos.castlingRights.CanCastle(c, QueenSide) &&
		pos.bbOccupied&bbQueenTravel == 0 {
		moves = append(moves, newMove(
			King.color(c),
			NoPiece,
			sqKing,
			sqQueenSide,
			NoSquare,
			NoPiece,
		))
	}

	return moves
}

func standardMoves(pos *Position) []Move {
	c := pos.turn
	bbAllowed := ^pos.getColor(c)
	bbPinned := pos.getPinned(c)

	var moves []Move
	for p1 := Pawn.color(c); p1 <= WhiteKing; p1 += 2 {
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

				p2 := pos.board.pieceByColor(s2, c.Other())
				if p1 == WhitePawn && s2.Rank() == Rank8 || p1 == BlackPawn && s2.Rank() == Rank1 {
					moves = append(moves,
						newMove(p1, p2, s1, s2, pos.enPassant, Queen.color(c)),
						newMove(p1, p2, s1, s2, pos.enPassant, Rook.color(c)),
						newMove(p1, p2, s1, s2, pos.enPassant, Bishop.color(c)),
						newMove(p1, p2, s1, s2, pos.enPassant, Knight.color(c)),
					)
				} else {
					moves = append(moves, newMove(p1, p2, s1, s2, pos.enPassant, NoPiece))
				}
			}
		}
	}
	return moves
}

func isCastleLegal(pos *Position, m Move) bool {
	var index int
	if m.HasTag(QueenSideCastle) {
		index |= 1
	}
	if pos.turn == Black {
		index |= 2
	}

	c := pos.turn.Other()
	check := castleChecks[index]

	if check.bbPawn&pos.getBitboard(Pawn.color(c)) > 0 ||
		check.bbKnight&pos.getBitboard(Knight.color(c)) > 0 ||
		check.bbKing&pos.getBitboard(King.color(c)) > 0 {
		return false
	}

	var bbBishop, bbRook bitboard
	for _, sq := range check.squares {
		bbBishop |= bishopAttacksBitboard(sq, pos.bbOccupied)
	}

	if bb := pos.getBitboard(Bishop.color(c)) | pos.getBitboard(Queen.color(c)); bbBishop&bb > 0 {
		return false
	}

	for _, sq := range check.squares {
		bbRook |= rookAttacksBitboard(sq, pos.bbOccupied)
	}

	return bbRook&(pos.getBitboard(Rook.color(c))|pos.getBitboard(Queen.color(c))) == 0
}

func movePinnedBitboard(sq Square, pos *Position, pt PieceType) bitboard {
	switch king := pos.getKingSquare(pos.turn); pt {
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
		pinner := pos.getPinner(pos.turn.Other())
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

func checkBitboard(sq Square, c Color, occupied, k, q, r, b, n, p bitboard) bitboard {
	bbRank := (q | r) & linearBitboard(sq, occupied, bbRanks[sq])
	bbFile := (q | r) & linearBitboard(sq, occupied, bbFiles[sq])
	bbDiag := (q | b) & linearBitboard(sq, occupied, bbDiagonals[sq])
	bbAnti := (q | b) & linearBitboard(sq, occupied, bbAntiDiagonals[sq])

	bbCheck := bbRank | bbFile | bbDiag | bbAnti
	bbCheck |= k & bbKingMoves[sq]
	bbCheck |= n & bbKnightMoves[sq]

	if c == Black {
		return bbCheck | p&bbBlackPawnCaptures[sq]
	}

	return bbCheck | p&bbWhitePawnCaptures[sq]
}

func pinnedBitboard(sq Square, occupied, blockers, queen, rook, bishop bitboard) (bitboard, bitboard) {
	var pinned, pinner bitboard

	if bbRanks[sq]&(queen|rook) > 0 {
		rPinner := xrayRankAttacksBitboard(sq, occupied, blockers) & (queen | rook)
		pinner |= rPinner

		for ; rPinner > 0; rPinner = rPinner.resetLSB() {
			s := rPinner.scanForward()
			pinned |= bbInBetween[sq][s] & blockers
		}
	}

	if bbFiles[sq]&(queen|rook) > 0 {
		rPinner := xrayFileAttacksBitboard(sq, occupied, blockers) & (queen | rook)
		pinner |= rPinner

		for ; rPinner > 0; rPinner = rPinner.resetLSB() {
			s := rPinner.scanForward()
			pinned |= bbInBetween[sq][s] & blockers
		}
	}

	if bbDiagonals[sq]&(queen|bishop) > 0 {
		bPinner := xrayDiagonalAttacksBitboard(sq, occupied, blockers) & (queen | bishop)
		pinner |= bPinner

		for ; bPinner > 0; bPinner = bPinner.resetLSB() {
			s := bPinner.scanForward()
			pinned |= bbInBetween[sq][s] & blockers
		}
	}

	if bbAntiDiagonals[sq]&(queen|bishop) > 0 {
		bPinner := xrayAntiDiagonalAttacksBitboard(sq, occupied, blockers) & (queen | bishop)
		pinner |= bPinner

		for ; bPinner > 0; bPinner = bPinner.resetLSB() {
			s := bPinner.scanForward()
			pinned |= bbInBetween[sq][s] & blockers
		}
	}

	return pinned, pinner
}

func xrayRankAttacksBitboard(sq Square, occupied, blockers bitboard) bitboard {
	attacks := linearBitboard(sq, occupied, bbRanks[sq])
	blockers &= attacks
	return attacks ^ linearBitboard(sq, occupied^blockers, bbRanks[sq])
}

func xrayFileAttacksBitboard(sq Square, occupied, blockers bitboard) bitboard {
	attacks := linearBitboard(sq, occupied, bbFiles[sq])
	blockers &= attacks
	return attacks ^ linearBitboard(sq, occupied^blockers, bbFiles[sq])
}

func xrayDiagonalAttacksBitboard(sq Square, occupied, blockers bitboard) bitboard {
	attacks := linearBitboard(sq, occupied, bbDiagonals[sq])
	blockers &= attacks
	return attacks ^ linearBitboard(sq, occupied^blockers, bbDiagonals[sq])
}

func xrayAntiDiagonalAttacksBitboard(sq Square, occupied, blockers bitboard) bitboard {
	attacks := linearBitboard(sq, occupied, bbAntiDiagonals[sq])
	blockers &= attacks
	return attacks ^ linearBitboard(sq, occupied^blockers, bbAntiDiagonals[sq])
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
