package search

import (
	"github.com/LucasSim0n/isabel/pkg/board"
	"github.com/LucasSim0n/isabel/pkg/cmn"
	"github.com/LucasSim0n/isabel/pkg/tt"
)

type Searcher struct {
	tt      *tt.TranspositionTable
	killers [64][2]cmn.Move
}

func NewSearcher() *Searcher {
	return &Searcher{
		tt: tt.NewTT(),
	}
}

func (s *Searcher) FindBestMove(b *board.Board, maxDepth int) cmn.Move {
	var move cmn.Move

	for depth := 1; depth <= maxDepth; depth++ {
		move = s.searchRoot(b, depth)
	}
	return move
}

func (s *Searcher) Search(b *board.Board, depth, alpha, beta, ply int) int {

	if entry, ok := s.tt.Probe(b.Hash, depth, alpha, beta); ok {
		return entry.Score
	}

	if depth == 0 {
		return s.Quiescence(b, alpha, beta, ply+1)
	}

	moves := b.GenerateLegalMoves()
	if len(moves) == 0 {
		inCheck := b.IsSquareAttacked(int(b.KingSq[b.SideToMove]), cmn.GetOpposite(b.SideToMove))
		if inCheck {
			return -(cmn.MateScore - ply)
		}

		return 0
	}

	var ttMove cmn.Move

	if entry, ok := s.tt.Get(b.Hash); ok {
		ttMove = entry.Move
	}

	s.sortMoves(moves, ttMove, ply)

	bestScore := -cmn.Infinity
	var bestMove cmn.Move
	alphaOrig := alpha

	for _, move := range moves {
		undo := b.MakeMove(move)

		score := -s.Search(
			b,
			depth-1,
			-beta,
			-alpha,
			ply+1,
		)

		b.UnmakeMove(move, undo)

		if score > bestScore {
			bestScore = score
			bestMove = move
		}

		if score > alpha {
			alpha = score
		}

		if alpha >= beta {

			if move.Flags&cmn.FlagCapture == 0 {
				s.killers[ply][1] = s.killers[ply][0]
				s.killers[ply][0] = move
			}

			break
		}
	}

	flag := tt.Exact
	switch {
	case bestScore <= alphaOrig:
		flag = tt.UpperBound
	case bestScore >= beta:
		flag = tt.LowerBound
	}

	s.tt.Store(b.Hash, depth, bestScore, flag, bestMove)

	return bestScore
}

func (s *Searcher) searchRoot(b *board.Board, depth int) cmn.Move {
	moves := b.GenerateLegalMoves()

	bestScore := -cmn.Infinity
	bestMove := moves[0]

	alpha := -cmn.Infinity
	beta := cmn.Infinity

	for _, move := range moves {
		undo := b.MakeMove(move)

		score := -s.Search(b, depth-1, -beta, -alpha, 0)

		b.UnmakeMove(move, undo)

		if score > bestScore {
			bestScore = score
			bestMove = move
		}

		if score > alpha {
			alpha = score
		}
	}

	return bestMove
}

func (s *Searcher) Quiescence(b *board.Board, alpha, beta, ply int) int {

	standPat := Evaluate(b)

	if standPat >= beta {
		return beta
	}

	if standPat > alpha {
		alpha = standPat
	}

	captures := b.GenerateLegalCaptures()
	if len(captures) == 0 {
		inCheck := b.IsSquareAttacked(int(b.KingSq[b.SideToMove]), cmn.GetOpposite(b.SideToMove))
		if inCheck {
			return -(cmn.MateScore - ply)
		}
	}

	s.sortMoves(captures, cmn.Move{}, ply)

	for _, move := range captures {
		undo := b.MakeMove(move)

		score := -s.Quiescence(b, -beta, -alpha, ply+1)

		b.UnmakeMove(move, undo)

		if score >= beta {
			return beta
		}

		if score > alpha {
			alpha = score
		}
	}

	return alpha
}
