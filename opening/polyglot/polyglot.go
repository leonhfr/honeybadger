// Package polyglot reads and provides an interface to polyglot opening book
// files (.bin) and provides an interface to.
package polyglot

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"sort"

	"github.com/notnil/chess"
)

// Book holds the opening book data.
type Book struct {
	positions map[uint64][]entry
}

type entry struct {
	move   uint16
	weight uint16
	learn  uint32
}

// WeightedMove is a single weight move. The weight is positive but no
// bounds are enforced.
type WeightedMove struct {
	Move   *chess.Move
	Weight int
}

// New returns a new empty Book.
func New() *Book {
	return &Book{
		positions: make(map[uint64][]entry),
	}
}

// Init takes a reader to binary data from Polyglot files (.bin) and
// initializes it. It may be called several times with data from
// different books and data will be merged. Weights will be returned as is
// so it is up to the user to ensures they are consistent across books.
func (b *Book) Init(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	scanner.Split(scanEntries)

	for scanner.Scan() {
		data := scanner.Bytes()
		key := binary.BigEndian.Uint64(data[0:8])
		move := binary.BigEndian.Uint16(data[8:10])
		weight := binary.BigEndian.Uint16(data[10:12])
		learn := binary.BigEndian.Uint32(data[12:16])

		entries, ok := b.positions[key]
		if ok {
			b.positions[key] = append(entries, entry{move, weight, learn})
		} else {
			b.positions[key] = []entry{{move, weight, learn}}
		}
	}

	return scanner.Err()
}

// Lookup takes a position and returns a sorted list of weighted moves.
// If the position is not found, nil is returned.
func (b *Book) Lookup(position *chess.Position) []WeightedMove {
	key := keyHash(position)
	entries, ok := b.positions[key]
	if !ok {
		return nil
	}

	var weightedMoves []WeightedMove
	validMoves := position.ValidMoves()

	for _, entry := range entries {
		from, to, promotion := parseMove(entry.move)
		fromPiece := position.Board().Piece(from)
		to = castlingDestination(from, to, fromPiece)

		for _, validMove := range validMoves {
			if from == validMove.S1() && to == validMove.S2() && promotion == validMove.Promo() {
				weightedMoves = append(weightedMoves, WeightedMove{validMove, int(entry.weight)})
				break
			}
		}
	}

	sort.Slice(weightedMoves, func(i, j int) bool {
		return weightedMoves[i].Weight > weightedMoves[j].Weight
	})

	return weightedMoves
}

// parseMove parses a raw move and returns origin and destination squares as
// well as promotion.
//
// A raw move is a bit field with the following meaning (bit 0 is the least
// significant bit):
//
//	===================================
//	0,1,2               to file
//	3,4,5               to row
//	6,7,8               from file
//	9,10,11             from row
//	12,13,14            promotion piece
func parseMove(move uint16) (from, to chess.Square, promotion chess.PieceType) {
	toFile := move & 7
	toRank := (move >> 3) & 7
	fromFile := (move >> 6) & 7
	fromRank := (move >> 9) & 7
	promo := (move >> 12) & 7

	from = chess.Square(8*fromRank + fromFile)
	to = chess.Square(8*toRank + toFile)
	promotion = promotionCodes[int(promo)]
	return
}

// promotionCodes is an array of piece types indexed by promotion codes,
// determined as follow:
//
//	none       0
//	knight     1
//	bishop     2
//	rook       3
//	queen      4
var promotionCodes = []chess.PieceType{
	chess.NoPieceType,
	chess.Knight,
	chess.Bishop,
	chess.Rook,
	chess.Queen,
}

// castlingDestination returns a new destination square if the move is a castling.
//
// Castling moves are unconventionally represented as follow:
//
//	white short      e1h1
//	white long       e1a1
//	black short      e8h8
//	black long       e8a8
//
// If the move is not castling, the original destination square is returned.
func castlingDestination(from, to chess.Square, fromPiece chess.Piece) chess.Square {
	switch {
	case from == chess.E1 && to == chess.H1 && fromPiece == chess.WhiteKing:
		return chess.G1
	case from == chess.E1 && to == chess.A1 && fromPiece == chess.WhiteKing:
		return chess.C1
	case from == chess.E8 && to == chess.H8 && fromPiece == chess.WhiteKing:
		return chess.G8
	case from == chess.E8 && to == chess.A8 && fromPiece == chess.WhiteKing:
		return chess.C8
	default:
		return to
	}
}

// scanEntries is a bufio.Scanner split function
//
// It returns data by chunks of 16 bytes.
func scanEntries(data []byte, atEOF bool) (advance int, token []byte, err error) {
	switch {
	case atEOF && len(data) == 0:
		// expected EOF
		return 0, nil, nil
	case len(data) >= 16:
		// return raw entry bytes
		return 16, data[0:16], nil
	case atEOF:
		// if at EOF and we still have data, the data is malformed
		return len(data), data, errors.New("expected data to be multiple of 16 bytes")
	default:
		// request more data
		return 0, nil, nil
	}
}
