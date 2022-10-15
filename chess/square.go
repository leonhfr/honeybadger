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
)

// File returns the square's file.
func (sq Square) File() File {
	return File(sq & 7)
}

// Rank returns the square's rank.
func (sq Square) Rank() Rank {
	return Rank(sq & 56)
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
	// FileA is the file A.
	FileA File = iota
	// FileB is the file B.
	FileB
	// FileC is the file C.
	FileC
	// FileD is the file D.
	FileD
	// FileE is the file E.
	FileE
	// FileF is the file F.
	FileF
	// FileG is the file G.
	FileG
	// FileH is the file H.
	FileH
)

func (f File) String() string {
	return fileChars[f : f+1]
}

// Rank is the rank of a square.
type Rank uint8

const (
	// Rank1 is the rank 1.
	Rank1 Rank = 8 * iota
	// Rank2 is the rank 2.
	Rank2
	// Rank3 is the rank 3.
	Rank3
	// Rank4 is the rank 4.
	Rank4
	// Rank5 is the rank 5.
	Rank5
	// Rank6 is the rank 6.
	Rank6
	// Rank7 is the rank 7.
	Rank7
	// Rank8 is the rank 8.
	Rank8
)

func (r Rank) String() string {
	return rankChars[r/8 : (r/8)+1]
}
