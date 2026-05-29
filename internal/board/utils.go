package board

import "fmt"

type Bitboard uint64

type Color uint8

const (
	White Color = iota
	Black
)

type PieceType uint8

const (
	Pawn PieceType = iota
	Knight
	Bishop
	Rook
	Queen
	King
)

var pieceValues = [6]int{
	100, // Pawn
	320, // Knight
	330, // Bishop
	500, // Rook
	900, // Queen
	0,   // King
}

const Infinity = 1000000
const MateScore = 900000

var pieceCharMap = map[rune]struct {
	Type  PieceType
	Color Color
}{
	'P': {Pawn, White}, 'N': {Knight, White}, 'B': {Bishop, White},
	'R': {Rook, White}, 'Q': {Queen, White}, 'K': {King, White},
	'p': {Pawn, Black}, 'n': {Knight, Black}, 'b': {Bishop, Black},
	'r': {Rook, Black}, 'q': {Queen, Black}, 'k': {King, Black},
}

const (
	rank1 Bitboard = 0x00000000000000FF
	rank2 Bitboard = 0x000000000000FF00
	rank7 Bitboard = 0x00FF000000000000
	rank8 Bitboard = 0xFF00000000000000

	notAFile Bitboard = 0xfefefefefefefefe
	notHFile Bitboard = 0x7f7f7f7f7f7f7f7f
)

func notationToSquare(notation string) (int, error) {
	if len(notation) != 2 {
		return 0, fmt.Errorf("invalid square notation: %s", notation)
	}

	fileChar := notation[0]
	rankChar := notation[1]

	file := int(fileChar - 'a')
	rank := int(rankChar - '1')

	if file < 0 || file > 7 || rank < 0 || rank > 7 {
		return 0, fmt.Errorf("square notation out of bounds: %s", notation)
	}

	return rank*8 + file, nil
}

func squareToNotation(sq int) string {
	if sq < 0 || sq > 63 {
		return "-"
	}
	file := rune('a' + (sq % 8))
	rank := rune('1' + (sq / 8))
	return string([]rune{file, rank})
}

func getOpposite(c Color) Color {
	if c == White {
		return Black
	}
	return White
}
