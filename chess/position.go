package chess

import (
	"fmt"
	"strings"
)

// CastlingRight represents the castling right of one combination of side and color.
type CastlingRight uint8

const (
	// CastleWhiteKing represents white's king castle.
	CastleWhiteKing CastlingRight = 1 << iota
	// CastleWhiteQueen represents white's queen castle.
	CastleWhiteQueen
	// CastleBlackKing represents black's king castle.
	CastleBlackKing
	// CastleBlackQueen represents black's queen castle.
	CastleBlackQueen
)

// Position represents the state of the game.
type Position struct {
	board           *board
	turn            Color
	castlingRights  CastlingRight
	enPassantSquare Square
	halfMoveClock   int
	fullMoves       int
}

// FromFEN creates a Position from a FEN string.
func FromFEN(fen string) (*Position, error) {
	p := &Position{}
	err := p.Unmarshal(fen)
	return p, err
}

// Unmarshal assumes text is in Forsyth–Edwards Notation.
func (p *Position) Unmarshal(fen string) error {
	fields := strings.Fields(strings.TrimSpace(fen))
	if len(fields) != 6 {
		return fmt.Errorf("invalid fen (%s), must have 6 fields", fen)
	}

	var err error

	p.board, err = fenBoard(fields[0])
	if err != nil {
		return err
	}

	p.turn, err = fenTurn(fields[1])
	if err != nil {
		return err
	}

	p.castlingRights, err = fenCastlingRights(fields[2])
	if err != nil {
		return err
	}

	p.enPassantSquare, err = fenEnPassantSquare(fields[3])
	if err != nil {
		return err
	}

	p.halfMoveClock, err = fenHalfMoveClock(fields[4])
	if err != nil {
		return err
	}

	p.fullMoves, err = fenFullMoves(fields[5])
	if err != nil {
		return err
	}

	return nil
}
