package chess

// Square is one of the 64 squares on a chess board.
type Square uint8

// NewSquare creates a new Square from a File and a Rank
func NewSquare(f File, r Rank) Square {
	return Square(f) + Square(r)
}

//nolint:revive
const (
	A1, B1, C1, D1, E1, F1, G1, H1 Square = 8*iota + 0, 8*iota + 1, 8*iota + 2,
		8*iota + 3, 8*iota + 4, 8*iota + 5, 8*iota + 6, 8*iota + 7
	A2, B2, C2, D2, E2, F2, G2, H2
	A3, B3, C3, D3, E3, F3, G3, H3
	A4, B4, C4, D4, E4, F4, G4, H4
	A5, B5, C5, D5, E5, F5, G5, H5
	A6, B6, C6, D6, E6, F6, G6, H6
	A7, B7, C7, D7, E7, F7, G7, H7
	A8, B8, C8, D8, E8, F8, G8, H8
	NoSquare Square = 64
)

// File returns the square's file.
func (sq Square) File() File {
	return File(sq & 7)
}

// Rank returns the square's rank.
func (sq Square) Rank() Rank {
	return Rank(sq & 56)
}

func (sq Square) bitboard() bitboard {
	return 1 << sq
}

func (sq Square) String() string {
	return sq.File().String() + sq.Rank().String()
}

const (
	fileChars = "abcdefgh"
	rankChars = "12345678"
)

// A File is the file of a square.
type File uint8

const (
	FileA File = iota // FileA is the file A.
	FileB             // FileB is the file B.
	FileC             // FileC is the file C.
	FileD             // FileD is the file D.
	FileE             // FileE is the file E.
	FileF             // FileF is the file F.
	FileG             // FileG is the file G.
	FileH             // FileH is the file H.
)

func (f File) String() string {
	return fileChars[f : f+1]
}

// Rank is the rank of a square.
type Rank uint8

const (
	Rank1 Rank = 8 * iota // Rank1 is the rank 1.
	Rank2                 // Rank2 is the rank 2.
	Rank3                 // Rank3 is the rank 3.
	Rank4                 // Rank4 is the rank 4.
	Rank5                 // Rank5 is the rank 5.
	Rank6                 // Rank6 is the rank 6.
	Rank7                 // Rank7 is the rank 7.
	Rank8                 // Rank8 is the rank 8.
)

func (r Rank) String() string {
	return rankChars[r/8 : (r/8)+1]
}

var strToSquareMap = map[string]Square{
	"a1": A1, "a2": A2, "a3": A3, "a4": A4, "a5": A5, "a6": A6, "a7": A7, "a8": A8,
	"b1": B1, "b2": B2, "b3": B3, "b4": B4, "b5": B5, "b6": B6, "b7": B7, "b8": B8,
	"c1": C1, "c2": C2, "c3": C3, "c4": C4, "c5": C5, "c6": C6, "c7": C7, "c8": C8,
	"d1": D1, "d2": D2, "d3": D3, "d4": D4, "d5": D5, "d6": D6, "d7": D7, "d8": D8,
	"e1": E1, "e2": E2, "e3": E3, "e4": E4, "e5": E5, "e6": E6, "e7": E7, "e8": E8,
	"f1": F1, "f2": F2, "f3": F3, "f4": F4, "f5": F5, "f6": F6, "f7": F7, "f8": F8,
	"g1": G1, "g2": G2, "g3": G3, "g4": G4, "g5": G5, "g6": G6, "g7": G7, "g8": G8,
	"h1": H1, "h2": H2, "h3": H3, "h4": H4, "h5": H5, "h6": H6, "h7": H7, "h8": H8,
}
