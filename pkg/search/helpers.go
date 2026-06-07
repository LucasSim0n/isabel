package search

import (
	"sort"

	"github.com/LucasSim0n/isabel/pkg/cmn"
)

func (s *Searcher) scoreMove(move, ttMove cmn.Move, ply int) int {
	if move == ttMove {
		return 1_000_000
	}

	if move.Flags&cmn.FlagCapture != 0 {
		return 100_000 + 10*PieceValues[move.Capture] - PieceValues[move.Piece]
	}

	if move == s.killers[ply][0] {
		return 90_000
	}

	if move == s.killers[ply][1] {
		return 80_000
	}

	return 0
}

func (s *Searcher) sortMoves(moves []cmn.Move, ttMove cmn.Move, ply int) {
	sort.Slice(moves, func(i, j int) bool {
		return s.scoreMove(moves[i], ttMove, ply) > s.scoreMove(moves[j], ttMove, ply)
	})
}

func getPSTValue(p cmn.PieceType, sq int, color cmn.Color) int {
	if color == cmn.Black {
		sq = 63 - sq
	}

	switch p {
	case cmn.King:
		return KingMap[sq]
	case cmn.Queen:
		return QueenMap[sq]
	case cmn.Rook:
		return RookMap[sq]
	case cmn.Bishop:
		return BishopMap[sq]
	case cmn.Knight:
		return KnightMap[sq]
	case cmn.Pawn:
		return PawnMap[sq]
	}
	return 0
}
