package search

import (
	"math/bits"

	"github.com/LucasSim0n/isabel/pkg/board"
	"github.com/LucasSim0n/isabel/pkg/cmn"
)

var PieceValues = [6]int{
	100, // Pawn
	320, // Knight
	330, // Bishop
	500, // Rook
	900, // Queen
	0,   // King
}

func Evaluate(b *board.Board) int {
	score := 0

	for p := cmn.Pawn; p <= cmn.Queen; p++ {

		whiteCount := bits.OnesCount64(uint64(b.Pieces[p]))

		blackCount := bits.OnesCount64(uint64(b.Pieces[p+6]))

		score += whiteCount * PieceValues[p]
		score -= blackCount * PieceValues[p]
	}

	if b.SideToMove == cmn.Black {
		score = -score
	}

	return score
}
