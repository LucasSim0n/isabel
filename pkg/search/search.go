package search

import (
	"errors"

	"github.com/LucasSim0n/isabel/pkg/board"
	"github.com/LucasSim0n/isabel/pkg/cmn"
	"github.com/LucasSim0n/isabel/pkg/tt"
)

type Searcher struct {
	tt *tt.TranspositionTable
}

func NewSearcher() *Searcher {
	return &Searcher{
		tt: tt.NewTT(),
	}
}

func (s *Searcher) Search(b *board.Board, depth, alpha, beta int) int {

	if entry, ok := s.tt.Probe(b.Hash, depth, alpha, beta); ok {
		return entry.Score
	}

	if depth == 0 {
		return s.Quiescence(b, alpha, beta)
	}

	moves := b.GenerateLegalMoves()

	if len(moves) == 0 {
		inCheck := b.IsSquareAttacked(int(b.KingSq[b.SideToMove]), cmn.GetOpposite(b.SideToMove))

		if inCheck {
			return -cmn.MateScore
		}

		return 0
	}

	var ttMove cmn.Move
	hasTTMove := false

	if entry, ok := s.tt.Get(b.Hash); ok {
		ttMove = entry.Move
		hasTTMove = true
	}

	if hasTTMove {
		for i, move := range moves {
			if ttMove == move {
				moves[0], moves[i] = moves[i], moves[0]
				break
			}
		}
	}

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

func (s *Searcher) FindBestMove(b *board.Board, depth int) (cmn.Move, error) {
	moves := b.GenerateLegalMoves()

	if len(moves) == 0 {
		return cmn.Move{}, errors.New("No legal moves found")
	}

	bestScore := -cmn.Infinity
	bestMove := moves[0]

	alpha := -cmn.Infinity
	beta := cmn.Infinity

	for _, move := range moves {

		undo := b.MakeMove(move)

		score := -s.Search(
			b,
			depth-1,
			-beta,
			-alpha,
		)

		b.UnmakeMove(move, undo)

		if score > bestScore {
			bestScore = score
			bestMove = move
		}

		if score > alpha {
			alpha = score
		}
	}

	return bestMove, nil
}

func (s *Searcher) Quiescence(b *board.Board, alpha, beta int) int {

	standPat := Evaluate(b)

	if standPat >= beta {
		return beta
	}

	if standPat > alpha {
		alpha = standPat
	}

	captures := b.GenerateLegalCaptures()

	for _, move := range captures {
		undo := b.MakeMove(move)

		score := -s.Quiescence(b, -beta, -alpha)

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
