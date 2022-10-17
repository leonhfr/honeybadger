package chess

import (
	"fmt"
	"strings"
)

// Position represents the state of the game.
type Position struct {
	board           *board
	turn            Color
	castlingRights  CastlingRights
	enPassantSquare Square
	halfMoveClock   int
	fullMoves       int
}

const startFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// StartingPosition returns the starting position.
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

func (pos *Position) String() string {
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
