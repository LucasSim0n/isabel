package board

import (
	"fmt"
	"math/bits"
	"strconv"
	"strings"
)

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
		s, err := notationToSquare(parts[3])
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

func (b *Board) GenerateMoves() *[]Move {
	var moves []Move

	b.generateKingMoves(&moves)
	b.generateKnightMoves(&moves)

	return &moves
}

func (b *Board) generateKingMoves(moves *[]Move) {
	color := b.SideToMove

	from := int(b.KingSq[color])
	attacks := kingAttacks[from]
	attacks &= ^b.Occupancy[color]

	b.iterateMoves(from, attacks, King, moves)

	b.generateCastlingMoves(moves)
}

func (b *Board) generateCastlingMoves(moves *[]Move) {
	color := b.SideToMove

	if color == White {
		// King side (e1 -> g1)
		if b.CastlingRights&0b1000 != 0 {
			if ((b.Occupancy[2]>>5)&1 == 0) && ((b.Occupancy[2]>>6)&1 == 0) {
				*moves = append(*moves, Move{
					From:  4,
					To:    6,
					Piece: King,
					Flags: FlagCastling,
				})
			}
		}

		// Queen side (e1 -> c1)
		if b.CastlingRights&0b0100 != 0 {
			if ((b.Occupancy[2]>>3)&1 == 0) &&
				((b.Occupancy[2]>>2)&1 == 0) &&
				((b.Occupancy[2]>>1)&1 == 0) {
				*moves = append(*moves, Move{
					From:  4,
					To:    2,
					Piece: King,
					Flags: FlagCastling,
				})
			}
		}
	} else {
		// Negro (simétrico)
		if b.CastlingRights&0b0010 != 0 {
			if ((b.Occupancy[2]>>61)&1 == 0) && ((b.Occupancy[2]>>62)&1 == 0) {
				*moves = append(*moves, Move{
					From:  60,
					To:    62,
					Piece: King,
					Flags: FlagCastling,
				})
			}
		}

		if b.CastlingRights&0b0001 != 0 {
			if ((b.Occupancy[2]>>59)&1 == 0) &&
				((b.Occupancy[2]>>58)&1 == 0) &&
				((b.Occupancy[2]>>57)&1 == 0) {
				*moves = append(*moves, Move{
					From:  60,
					To:    58,
					Piece: King,
					Flags: FlagCastling,
				})
			}
		}
	}
}

func (b *Board) generateKnightMoves(moves *[]Move) {
	color := b.SideToMove
	index := int(Knight)
	if color == Black {
		index += 6
	}

	bb := b.Pieces[index]
	for bb != 0 {
		from := bits.TrailingZeros64(uint64(bb))
		attacks := knightAttacks[from]
		attacks &= ^b.Occupancy[color]

		b.iterateMoves(from, attacks, Knight, moves)

		bb &= bb - 1
	}
}

func (b *Board) generatePawnMoves(moves *[]Move) {
	color := b.SideToMove
	enemy := getEnemy(color)
	forward := 8
	index := Pawn

	if color == Black {
		index += 6
		forward = -8
	}

	pawns := b.Pieces[index]

	empty := ^b.Occupancy[2]

	// =====================
	// 1. SIMPLE PUSH
	// =====================
	var singlePush Bitboard
	if color == White {
		singlePush = (pawns << 8) & empty
	} else {
		singlePush = (pawns >> 8) & empty
	}

	b.iteratePawnMoves(singlePush, forward, moves)

	// =====================
	// 2. DOUBLE PUSH
	// =====================
	if color == White {
		p := pawns & rank2
		s1 := (p << 8) & empty
		doublePush := (s1 << 8) & empty

		b.iteratePawnMoves(doublePush, forward*2, moves)
	} else {
		p := pawns & rank7
		s1 := (p >> 8) & empty
		doublePush := (s1 >> 8) & empty

		b.iteratePawnMoves(doublePush, forward*2, moves)
	}

	// =====================
	// 3. CAPTURES
	// =====================
	var left, right Bitboard

	if color == White {
		left = (pawns << 7) & notHFile & b.Occupancy[enemy]
		right = (pawns << 9) & notAFile & b.Occupancy[enemy]

		b.iteratePawnCaptures(left, 7, moves)
		b.iteratePawnCaptures(right, 9, moves)

	} else {
		left = (pawns >> 9) & notHFile & b.Occupancy[enemy]
		right = (pawns >> 7) & notAFile & b.Occupancy[enemy]

		b.iteratePawnCaptures(left, -7, moves)
		b.iteratePawnCaptures(right, -9, moves)
	}
}

func (b *Board) iteratePawnMoves(bb Bitboard, offset int, moves *[]Move) {
	for bb != 0 {
		to := bits.TrailingZeros64(uint64(bb))
		from := to - offset

		move := NewMove(from, to, Pawn)

		if offset == 16 || offset == -16 {
			move.Flags |= FlagDoublePawnPush
		}
		*moves = append(*moves, move)

		bb &= bb - 1
	}

}

func (b *Board) iteratePawnCaptures(bb Bitboard, offset int, moves *[]Move) {
	for bb != 0 {
		to := bits.TrailingZeros64(uint64(bb))
		from := to - offset

		move := NewMove(from, to, Pawn)
		move.Flags |= FlagCapture
		move.Capture = b.getPieceAt(to)

		*moves = append(*moves, move)

		bb &= bb - 1
	}
}

func (b *Board) iterateMoves(from int, attacks Bitboard, piece PieceType, moves *[]Move) {

	enemy := getEnemy(b.SideToMove)

	for attacks != 0 {
		to := bits.TrailingZeros64(uint64(attacks))

		move := NewMove(from, to, piece)

		if (b.Occupancy[enemy]>>to)&1 != 0 {
			move.Flags |= FlagCapture

			move.Capture = b.getPieceAt(to)
		}

		*moves = append(*moves, move)

		attacks &= attacks - 1
	}
}

func (b *Board) getPieceAt(sq int) PieceType {
	for i := range 12 {
		if (b.Pieces[i]>>sq)&1 != 0 {
			return PieceType(i % 6)
		}
	}
	return Pawn
}

/*  TODO:
// 	func (b *Board) ToFEN() string
// 	func (b *Board) MakeMove(move *Move) *Board
// 	func (b *Board) IsLegalMove(move *Move) bool
*/

func (b *Board) String() string {
	var sb strings.Builder
	pieceChars := []rune{'P', 'N', 'B', 'R', 'Q', 'K', 'p', 'n', 'b', 'r', 'q', 'k'}

	sb.WriteString("  +-----------------+\n")
	for rank := 7; rank >= 0; rank-- {
		fmt.Fprintf(&sb, "%d | ", rank+1)
		for file := range 8 {
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
		fmt.Fprintf(&sb, "En Passant square: %s\n", squareToNotation(b.EnPassant))
	} else {
		sb.WriteString("En Passant square: -\n")
	}
	fmt.Fprintf(&sb, "Halfmove clock: %d\n", b.HalfmoveClock)
	fmt.Fprintf(&sb, "Fullmove number: %d\n", b.FullmoveNumber)
	fmt.Fprintf(&sb, "White King: %s, Black King: %s\n", squareToNotation(int(b.KingSq[White])), squareToNotation(int(b.KingSq[Black])))

	return sb.String()
}
