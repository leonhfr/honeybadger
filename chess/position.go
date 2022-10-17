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
	p := &Position{}

	p.board, err = fenBoard(fields[0])
	if err != nil {
		return nil, err
	}

	p.turn, err = fenTurn(fields[1])
	if err != nil {
		return nil, err
	}

	p.castlingRights, err = fenCastlingRights(fields[2])
	if err != nil {
		return nil, err
	}

	p.enPassantSquare, err = fenEnPassantSquare(fields[3])
	if err != nil {
		return nil, err
	}

	p.halfMoveClock, err = fenHalfMoveClock(fields[4])
	if err != nil {
		return nil, err
	}

	p.fullMoves, err = fenFullMoves(fields[5])
	if err != nil {
		return nil, err
	}

	return p, nil
}

// SquareMap returns the map from square to pieces.
func (p Position) SquareMap() SquareMap {
	return p.board.squareMap()
}

// Piece returns the piece present in square sq. Returns NoPiece if there aren't any.
func (p Position) Piece(sq Square) Piece {
	return p.board.piece(sq)
}

// Turn returns the color of the next player to move in this position.
func (p Position) Turn() Color {
	return p.turn
}

// CastlingRights returns the castling rights of the position.
func (p Position) CastlingRights() CastlingRights {
	return p.castlingRights
}

// EnPassantSquare returns the en passant square.
func (p Position) EnPassantSquare() Square {
	return p.enPassantSquare
}

// HalfMoveClock returns the half-move clock.
func (p Position) HalfMoveClock() int {
	return p.halfMoveClock
}

// FullMoves returns the full moves count.
func (p Position) FullMoves() int {
	return p.fullMoves
}

func (p *Position) String() string {
	sq := "-"
	if p.enPassantSquare != NoSquare {
		sq = p.enPassantSquare.String()
	}

	return fmt.Sprintf(
		"%s %s %s %s %d %d",
		p.board.String(),
		p.turn.String(),
		p.castlingRights.String(),
		sq,
		p.halfMoveClock,
		p.fullMoves,
	)
}
