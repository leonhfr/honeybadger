package chess

import (
	"fmt"
	"strings"
)

// Position represents the state of the game.
type Position struct {
	board
	turn            Color
	castlingRights  CastlingRights
	enPassantSquare Square
	halfMoveClock   int
	fullMoves       int
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

	pos.enPassantSquare, err = fenEnPassantSquare(fields[3])
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

// Piece returns the piece present in square sq. Returns NoPiece if there aren't any.
func (pos Position) Piece(sq Square) Piece {
	return pos.board.piece(sq)
}

// Turn returns the color of the next player to move in this position.
func (pos Position) Turn() Color {
	return pos.turn
}

// CastlingRights returns the castling rights of the position.
func (pos Position) CastlingRights() CastlingRights {
	return pos.castlingRights
}

// EnPassantSquare returns the en passant square.
func (pos Position) EnPassantSquare() Square {
	return pos.enPassantSquare
}

// HalfMoveClock returns the half-move clock.
func (pos Position) HalfMoveClock() int {
	return pos.halfMoveClock
}

// FullMoves returns the full moves count.
func (pos Position) FullMoves() int {
	return pos.fullMoves
}

func (pos Position) String() string {
	sq := "-"
	if pos.enPassantSquare != NoSquare {
		sq = pos.enPassantSquare.String()
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
func (pos *Position) PseudoMoves() []Move {
	return pseudoMoves(pos)
}

// MakeMove plays a move on a position and checks whether it is valid.
func (pos Position) MakeMove(m Move) (*Position, bool) {
	board := pos.board.copy()
	board.makeMove(m)

	if !(m.HasTag(KingSideCastle) || m.HasTag(QueenSideCastle)) &&
		isInCheck(&Position{board: board, turn: pos.turn}) {
		return nil, false
	}

	halfMoveClock := pos.halfMoveClock
	if m.P1().Type() == Pawn || m.HasTag(Capture) {
		halfMoveClock = 0
	} else {
		halfMoveClock++
	}

	fullMoves := pos.fullMoves
	if pos.turn == Black {
		fullMoves++
	}

	return &Position{
		board:           board,
		turn:            pos.turn.Other(),
		castlingRights:  pos.moveCastlingRights(m),
		enPassantSquare: pos.moveEnPassantSquare(m),
		halfMoveClock:   halfMoveClock,
		fullMoves:       fullMoves,
	}, true
}

func (pos Position) moveCastlingRights(m Move) CastlingRights {
	switch p1, s1, s2 := m.P1(), m.S1(), m.S2(); {
	case p1 == WhiteKing:
		return pos.castlingRights & ^(CastleWhiteKing | CastleWhiteQueen)
	case p1 == BlackKing:
		return pos.castlingRights & ^(CastleBlackKing | CastleBlackQueen)
	case (p1 == WhiteRook && s1 == A1) || s2 == A1:
		return pos.castlingRights & ^CastleWhiteQueen
	case (p1 == WhiteRook && s1 == H1) || s2 == H1:
		return pos.castlingRights & ^CastleWhiteKing
	case (p1 == BlackRook && s1 == A8) || s2 == A8:
		return pos.castlingRights & ^CastleBlackQueen
	case (p1 == BlackRook && s1 == H8) || s2 == H8:
		return pos.castlingRights & ^CastleBlackKing
	default:
		return pos.castlingRights
	}
}

func (pos Position) moveEnPassantSquare(m Move) Square {
	if m.P1().Type() != Pawn {
		return NoSquare
	}

	switch s1, s2 := m.S1(), m.S2(); {
	case pos.turn == White &&
		s1.bitboard()&bbRank2 > 0 &&
		s2.bitboard()&bbRank4 > 0:
		return s2 - 8
	case pos.turn == Black &&
		s1.bitboard()&bbRank7 > 0 &&
		s2.bitboard()&bbRank5 > 0:
		return s2 + 8
	default:
		return NoSquare
	}
}
