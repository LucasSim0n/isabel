package search

import (
	"math/bits"

	"github.com/LucasSim0n/isabel/pkg/board"
	"github.com/LucasSim0n/isabel/pkg/cmn"
)

func Evaluate(b *board.Board) int {
	score := 0

	for p := cmn.Pawn; p <= cmn.King; p++ {

		whitePiecesBB := b.Pieces[p]
		blackPiecesBB := b.Pieces[p+6]

		whiteCount := bits.OnesCount64(uint64(whitePiecesBB))

		blackCount := bits.OnesCount64(uint64(blackPiecesBB))

		score += whiteCount * PieceValues[p]
		score -= blackCount * PieceValues[p]

		for sq := range 64 {
			if whitePiecesBB&(1<<sq) != 0 {
				score += getPSTValue(p, sq, cmn.White)
			}
			if blackPiecesBB&(1<<sq) != 0 {
				score -= getPSTValue(p, sq, cmn.Black)
			}
		}
	}

	if b.SideToMove == cmn.Black {
		score = -score
	}

	return score
}
