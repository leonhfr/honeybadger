// Package chess provides types and functions to handle chess positions.
package chess

func pseudoMoves(pos *Position) []Move {
	switch checks := pos.getCheck(pos.turn.Other()); checks.ones() {
	case 2: // double check
		return checkFlightMoves(pos)
	default: // no check & single check
		return append(standardMoves(pos), castlingMoves(pos)...)
	}
}

// king moves to non attacked squares
func checkFlightMoves(pos *Position) []Move {
	c, op := pos.turn, pos.turn.Other()

	bbFlight := bbKingMoves[pos.getKingSquare(c)]
	bbFlight &= ^pos.getColor(c)                    // possible moves
	bbFlight &= ^bbKingMoves[pos.getKingSquare(op)] // enemy king attacks

	// pawn attacks
	if c == White {
		bbFlight &= ^((pos.bbBlackPawn & ^bbFileA) >> 9)
		bbFlight &= ^((pos.bbBlackPawn & ^bbFileH) >> 7)
	} else {
		bbFlight &= ^((pos.bbWhitePawn &^ bbFileA) << 9)
		bbFlight &= ^((pos.bbWhitePawn &^ bbFileH) << 7)
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

	moves := make([]Move, 0, 8)
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

	moves := make([]Move, 0, 128)
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
