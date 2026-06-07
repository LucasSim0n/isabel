package search

import (
	"sort"

	"github.com/LucasSim0n/isabel/pkg/cmn"
)

func scoreMove(move, ttMove cmn.Move) int {
	if move == ttMove {
		return 1_000_000
	}

	if move.Flags&cmn.FlagCapture != 0 {
		return 100_000 + 10*PieceValues[move.Capture] - PieceValues[move.Piece]
	}

	return 0
}

func sortMoves(moves []cmn.Move, ttMove cmn.Move) {
	sort.Slice(moves, func(i, j int) bool {
		return scoreMove(moves[i], ttMove) > scoreMove(moves[j], ttMove)
	})
}
