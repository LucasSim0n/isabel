package board

import (
	"fmt"
	"strings"

	"github.com/LucasSim0n/isabel/pkg/cmn"
)

/*  TODO:
// 	func (b *Board) ToFEN() string
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
		if b.SideToMove == cmn.White {
			return "White"
		} else {
			return "Black"
		}
	}())

	fmt.Fprintf(&sb, "Castling rights: %04b\n", b.CastlingRights)
	if b.EnPassant != -1 {
		fmt.Fprintf(&sb, "En Passant square: %s\n", cmn.SquareToNotation(b.EnPassant))
	} else {
		sb.WriteString("En Passant square: -\n")
	}
	fmt.Fprintf(&sb, "Halfmove clock: %d\n", b.HalfmoveClock)
	fmt.Fprintf(&sb, "Fullmove number: %d\n", b.FullmoveNumber)
	fmt.Fprintf(&sb, "White cmn.King: %s, cmn.Black cmn.King: %s\n", cmn.SquareToNotation(int(b.KingSq[cmn.White])), cmn.SquareToNotation(int(b.KingSq[cmn.Black])))
	fmt.Fprintf(&sb, "White cmn.King: %d, cmn.Black cmn.King: %d\n", b.KingSq[cmn.White], b.KingSq[cmn.Black])

	return sb.String()
}
