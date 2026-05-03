package board

import (
	"fmt"
	"math/bits"
	"strconv"
	"strings"
)

type Bitboard uint64

type Color int

const (
	White Color = iota
	Black
)

type PieceType uint

const (
	Pawn PieceType = iota
	Knight
	Bishop
	Rook
	Queen
	King
)

var pieceCharMap = map[rune]struct {
	Type  PieceType
	Color Color
}{
	'P': {Pawn, White}, 'N': {Knight, White}, 'B': {Bishop, White},
	'R': {Rook, White}, 'Q': {Queen, White}, 'K': {King, White},
	'p': {Pawn, Black}, 'n': {Knight, Black}, 'b': {Bishop, Black},
	'r': {Rook, Black}, 'q': {Queen, Black}, 'k': {King, Black},
}

type Board struct {
	Pieces    [12]Bitboard
	Occupancy [3]Bitboard

	SideToMove     Color
	CastlingRights uint8
	EnPassant      int
	HalfmoveClock  int
	FullmoveNumber int

	KingSq [2]uint8

	Hash uint64
}

func NewBoard(fen string) (*Board, error) {
	b := &Board{
		EnPassant: -1,
	}

	parts := strings.Fields(fen)
	if len(parts) < 6 {
		return nil, fmt.Errorf("Incomplete fen: %s", fen)
	}

	piecePlacement := parts[0]
	rank := 7
	file := 0

	for _, char := range piecePlacement {
		if char == '/' {
			rank--
			file = 0
			continue
		}

		if num, err := strconv.Atoi(string(char)); err == nil {
			file += num

		} else {
			pieceInfo, ok := pieceCharMap[char]
			if !ok {
				return nil, fmt.Errorf("Unknown FEN char: %s", fen)
			}

			square := rank*8 + file

			index := int(pieceInfo.Type)
			if pieceInfo.Color == Black {
				index += 6
			}

			b.Pieces[index] |= (1 << square)

			file++
		}
	}

	for i := range 6 {
		b.Occupancy[White] |= b.Pieces[i]
		b.Occupancy[Black] |= b.Pieces[i+6]
	}
	b.Occupancy[2] = b.Occupancy[White] | b.Occupancy[Black]

	whiteKingBB := b.Pieces[King]
	b.KingSq[White] = uint8(bits.TrailingZeros64(uint64(whiteKingBB)))

	blackKingBB := b.Pieces[King+6]
	b.KingSq[Black] = uint8(bits.TrailingZeros64(uint64(blackKingBB)))

	b.SideToMove = White
	if parts[1] == "b" {
		b.SideToMove = Black
	}

	b.CastlingRights = 0
	for _, c := range parts[2] {
		switch c {
		case 'K':
			b.CastlingRights |= 0b1000
		case 'Q':
			b.CastlingRights |= 0b0100
		case 'k':
			b.CastlingRights |= 0b0010
		case 'q':
			b.CastlingRights |= 0b0001
		}
	}

	if parts[3] != "-" {
		s, err := NotationToSquare(parts[3])
		if err != nil {
			return nil, fmt.Errorf("Incorrect fen: %s", fen)
		}
		b.EnPassant = s
	}

	halfmove, err := strconv.Atoi(parts[4])
	if err != nil {
		return nil, fmt.Errorf("Incorrect fen: %s", fen)
	}
	b.HalfmoveClock = halfmove

	fullmove, err := strconv.Atoi(parts[5])
	if err != nil {
		return nil, fmt.Errorf("Incorrect fen: %s", fen)
	}
	b.FullmoveNumber = fullmove

	return b, nil
}

func NotationToSquare(notation string) (int, error) {
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

func SquareToNotation(sq int) string {
	if sq < 0 || sq > 63 {
		return "-"
	}
	file := rune('a' + (sq % 8))
	rank := rune('1' + (sq / 8))
	return string([]rune{file, rank})
}

// func (b *Board) ToFEN() string
//
// func (b *Board) MakeMove(move *Move) *Board
//
// func (b *Board) IsLegalMove(move *Move) bool

func (b *Board) String() string {
	var sb strings.Builder
	pieceChars := []rune{'P', 'N', 'B', 'R', 'Q', 'K', 'p', 'n', 'b', 'r', 'q', 'k'}

	sb.WriteString("  +-----------------+\n")
	for rank := 7; rank >= 0; rank-- { // Desde la fila 8 hasta la 1
		fmt.Fprintf(&sb, "%d | ", rank+1)
		for file := range 8 { // Desde la columna A hasta la H
			square := rank*8 + file
			found := false
			for pieceTypeIndex := range 12 {
				if (b.Pieces[pieceTypeIndex]>>square)&1 != 0 {
					sb.WriteRune(pieceChars[pieceTypeIndex])
					found = true
					break
				}
			}
			if !found {
				sb.WriteRune('.')
			}
			sb.WriteRune(' ')
		}
		sb.WriteString("|\n")
	}
	sb.WriteString("  +-----------------+\n")
	sb.WriteString("    a b c d e f g h\n")

	fmt.Fprintf(&sb, "Side to move: %s\n", func() string {
		if b.SideToMove == White {
			return "White"
		} else {
			return "Black"
		}
	}())

	fmt.Fprintf(&sb, "Castling rights: %04b\n", b.CastlingRights)
	if b.EnPassant != -1 {
		fmt.Fprintf(&sb, "En Passant square: %s\n", SquareToNotation(b.EnPassant)) // Necesitas SquareToNotation
	} else {
		sb.WriteString("En Passant square: -\n")
	}
	fmt.Fprintf(&sb, "Halfmove clock: %d\n", b.HalfmoveClock)
	fmt.Fprintf(&sb, "Fullmove number: %d\n", b.FullmoveNumber)
	fmt.Fprintf(&sb, "White King: %s, Black King: %s\n", SquareToNotation(int(b.KingSq[White])), SquareToNotation(int(b.KingSq[Black])))

	return sb.String()
}
