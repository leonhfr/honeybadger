package chess

import (
	"fmt"
	"strconv"
	"strings"
)

func fenBoard(field string) (*board, error) {
	rankFields := strings.Split(field, "/")
	if len(rankFields) != 8 {
		return nil, fmt.Errorf("invalid fen board (%s)", field)
	}

	m := SquareMap{}
	for i, rankField := range rankFields {
		fileMap, err := fenFileMap(rankField)
		if err != nil {
			return nil, err
		}
		rank := Rank(56 - 8*i)
		for f, p := range fileMap {
			m[NewSquare(f, rank)] = p
		}
	}

	return newBoard(m), nil
}

func fenFileMap(rankField string) (map[File]Piece, error) {
	m := map[File]Piece{}
	file := FileA
	for _, r := range rankField {
		if p, ok := fenPieceMap[r]; ok {
			m[file] = p
			file++
		} else if '1' <= r && r <= '8' {
			file += File(r - '0')
		} else {
			return nil, fmt.Errorf("invalid fen rank field (%s)", rankField)
		}
	}

	if file != FileH+1 {
		return nil, fmt.Errorf("invalid fen rank field (%s)", rankField)
	}
	return m, nil
}

func fenTurn(field string) (Color, error) {
	turn, ok := fenTurnMap[field]
	if !ok {
		return White, fmt.Errorf("invalid fen turn (%s)", field)
	}
	return turn, nil
}

func fenCastlingRights(field string) (CastlingRights, error) {
	for _, s := range []string{"K", "Q", "k", "q", "-"} {
		if strings.Count(field, s) > 1 {
			return 0, fmt.Errorf("invalid fen castling rights (%s)", field)
		}
	}
	var castlingRights CastlingRights
	for _, r := range field {
		switch r {
		case 'K':
			castlingRights |= CastleWhiteKing
		case 'Q':
			castlingRights |= CastleWhiteQueen
		case 'k':
			castlingRights |= CastleBlackKing
		case 'q':
			castlingRights |= CastleBlackQueen
		case '-':
		default:
			return 0, fmt.Errorf("invalid fen castling rights (%s)", field)
		}
	}
	return castlingRights, nil
}

func fenEnPassantSquare(field string) (Square, error) {
	if field == "-" {
		return NoSquare, nil
	}
	sq, ok := strToSquareMap[field]
	if !ok || !(sq.Rank() == Rank3 || sq.Rank() == Rank6) {
		return NoSquare, fmt.Errorf("invalid fen en passant square (%s)", field)
	}
	return sq, nil
}

func fenHalfMoveClock(field string) (int, error) {
	halfMoveClock, err := strconv.Atoi(field)
	if err != nil || halfMoveClock < 0 {
		return 0, fmt.Errorf("invalid fen full moves count (%s)", field)
	}
	return halfMoveClock, nil
}

func fenFullMoves(field string) (int, error) {
	fullMoves, err := strconv.Atoi(field)
	if err != nil || fullMoves < 1 {
		return 0, fmt.Errorf("invalid fen full moves count (%s)", field)
	}
	return fullMoves, nil
}

var (
	fenPieceMap = map[rune]Piece{
		'K': WhiteKing,
		'Q': WhiteQueen,
		'R': WhiteRook,
		'B': WhiteBishop,
		'N': WhiteKnight,
		'P': WhitePawn,
		'k': BlackKing,
		'q': BlackQueen,
		'r': BlackRook,
		'b': BlackBishop,
		'n': BlackKnight,
		'p': BlackPawn,
	}

	fenTurnMap = map[string]Color{
		"w": White,
		"b": Black,
	}
)
