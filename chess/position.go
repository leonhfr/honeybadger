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

func (cr CastlingRight) String() string {
	var rights string
	if (cr & CastleWhiteKing) > 0 {
		rights += "K"
	}
	if (cr & CastleWhiteQueen) > 0 {
		rights += "Q"
	}
	if (cr & CastleBlackKing) > 0 {
		rights += "k"
	}
	if (cr & CastleBlackQueen) > 0 {
		rights += "q"
	}
	return rights
}

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
