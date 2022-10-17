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
