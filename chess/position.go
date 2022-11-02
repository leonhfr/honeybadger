package chess

import (
	"fmt"
	"strings"
)

// Position represents the state of the game.
type Position struct {
	board
	turn           Color
	castlingRights CastlingRights
	enPassant      Square
	halfMoveClock  uint8
	fullMoves      uint8
}

const startFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// StartingPosition returns the starting position.
//
// FEN: rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
func StartingPosition() *Position {
	p, _ := FromFEN(startFEN)
	return p
}

// FromFEN creates a Position from a FEN string.
func FromFEN(fen string) (*Position, error) {
	fields := strings.Fields(strings.TrimSpace(fen))
	if len(fields) != 6 {
		return nil, fmt.Errorf("invalid fen (%s), must have 6 fields", fen)
	}

	var err error
	pos := &Position{}

	pos.board, err = fenBoard(fields[0])
	if err != nil {
		return nil, err
	}

	pos.turn, err = fenTurn(fields[1])
	if err != nil {
		return nil, err
	}

	pos.castlingRights, err = fenCastlingRights(fields[2])
	if err != nil {
		return nil, err
	}

	pos.enPassant, err = fenEnPassantSquare(fields[3])
	if err != nil {
		return nil, err
	}

	pos.halfMoveClock, err = fenHalfMoveClock(fields[4])
	if err != nil {
		return nil, err
	}

	pos.fullMoves, err = fenFullMoves(fields[5])
	if err != nil {
		return nil, err
	}

	return pos, nil
}

// SquareMap returns the map from square to pieces.
func (pos Position) SquareMap() SquareMap {
	return pos.board.squareMap()
}

// PieceAt returns the piece present in square sq. Returns NoPiece if there aren't any.
func (pos Position) PieceAt(sq Square) Piece {
	return pos.board.pieceAt(sq)
}

// Turn returns the color of the next player to move in this position.
func (pos Position) Turn() Color {
	return pos.turn
}

// CastlingRights returns the castling rights of the position.
func (pos Position) CastlingRights() CastlingRights {
	return pos.castlingRights
}

// EnPassant returns the en passant square.
func (pos Position) EnPassant() Square {
	return pos.enPassant
}

// HalfMoveClock returns the half-move clock.
func (pos Position) HalfMoveClock() uint8 {
	return pos.halfMoveClock
}

// FullMoves returns the full moves count.
func (pos Position) FullMoves() uint8 {
	return pos.fullMoves
}

// InCheck indicates whether the position is in check.
func (pos Position) InCheck() bool {
	return pos.getCheck(pos.turn.Other()) > 0
}

func (pos Position) String() string {
	sq := "-"
	if pos.enPassant != NoSquare {
		sq = pos.enPassant.String()
	}

	return fmt.Sprintf(
		"%s %s %s %s %d %d",
		pos.board.String(),
		pos.turn.String(),
		pos.castlingRights.String(),
		sq,
		pos.halfMoveClock,
		pos.fullMoves,
	)
}

// PseudoMoves returns a list of pseudo moves.
//
// If the current player is not in check, some of the moves may not be legal.
// If the current player is in check, all of the moves are guaranteed to be legal.
func (pos *Position) PseudoMoves() []Move {
	return pseudoMoves(pos)
}

// MakeMove makes a move on a position and checks whether it is valid.
// Returns metadata that can be used to unmake the move and a boolean
// indicating the validity of the move.
func (pos *Position) MakeMove(m Move) (Metadata, bool) {
	metadata := newMetadata(pos.turn, pos.castlingRights,
		pos.halfMoveClock, pos.fullMoves, pos.enPassant)

	if (m.HasTag(KingSideCastle) || m.HasTag(QueenSideCastle)) && !isCastleLegal(pos, m) {
		return metadata, false
	}

	if pos.board.makeMoveBoard(m); pos.getCheck(pos.turn.Other()) > 0 {
		pos.board.unmakeMoveBoard(m)

		return metadata, false
	}

	pos.turn = pos.turn.Other()
	pos.castlingRights = moveCastlingRights(pos.castlingRights, m)
	pos.enPassant = moveEnPassantMove(m)

	if m.P1().Type() == Pawn || m.HasTag(Capture) {
		pos.halfMoveClock = 0
	} else {
		pos.halfMoveClock++
	}

	if pos.turn == White {
		pos.fullMoves++
	}

	return metadata, true
}

// UnmakeMove unmakes a move and restores the previous position.
func (pos *Position) UnmakeMove(m Move, meta Metadata) {
	pos.board.unmakeMoveBoard(m)
	pos.turn = meta.turn()
	pos.castlingRights = meta.castleRights()
	pos.enPassant = meta.enPassant()
	pos.halfMoveClock = meta.halfMoveClock()
	pos.fullMoves = meta.fullMoves()
}

// Hash returns a Zobrist hash of the position.
//
// The hash is compatible with polyglot files.
func (pos *Position) Hash() uint64 {
	return zobristHash(pos)
}

// Copy returns a copy of the position.
func (pos *Position) Copy() *Position {
	return &Position{
		board:          pos.board.copyBoard(),
		turn:           pos.turn,
		castlingRights: pos.castlingRights,
		enPassant:      pos.enPassant,
		halfMoveClock:  pos.halfMoveClock,
		fullMoves:      pos.fullMoves,
	}
}

func moveCastlingRights(cr CastlingRights, m Move) CastlingRights {
	switch p1, s1, s2 := m.P1(), m.S1(), m.S2(); {
	case p1 == WhiteKing:
		return cr & ^(CastleWhiteKing | CastleWhiteQueen)
	case p1 == BlackKing:
		return cr & ^(CastleBlackKing | CastleBlackQueen)
	case (p1 == WhiteRook && s1 == A1) || s2 == A1:
		return cr & ^CastleWhiteQueen
	case (p1 == WhiteRook && s1 == H1) || s2 == H1:
		return cr & ^CastleWhiteKing
	case (p1 == BlackRook && s1 == A8) || s2 == A8:
		return cr & ^CastleBlackQueen
	case (p1 == BlackRook && s1 == H8) || s2 == H8:
		return cr & ^CastleBlackKing
	default:
		return cr
	}
}

func moveEnPassantMove(m Move) Square {
	if m.P1().Type() != Pawn {
		return NoSquare
	}

	switch c, s1, s2 := m.P1().Color(), m.S1(), m.S2(); {
	case c == White &&
		s1.bitboard()&bbRank2 > 0 &&
		s2.bitboard()&bbRank4 > 0:
		return s2 - 8
	case c == Black &&
		s1.bitboard()&bbRank7 > 0 &&
		s2.bitboard()&bbRank5 > 0:
		return s2 + 8
	default:
		return NoSquare
	}
}
