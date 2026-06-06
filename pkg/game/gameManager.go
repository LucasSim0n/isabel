package game

import (
	"github.com/LucasSim0n/isabel/pkg/board"
	"github.com/LucasSim0n/isabel/pkg/cmn"
	"github.com/LucasSim0n/isabel/pkg/search"
)

type GameManager struct {
	board    *board.Board
	searcher *search.Searcher
}

func NewGameManager() *GameManager {
	return &GameManager{
		board:    board.NewStartPosGame(),
		searcher: search.NewSearcher(),
	}
}

func (g *GameManager) parseMove(moveStr string) (cmn.Move, error) {

	from, err := cmn.NotationToSquare(moveStr[0:2])
	if err != nil {
		return cmn.Move{}, err
	}

	to, err := cmn.NotationToSquare(moveStr[2:4])
	if err != nil {
		return cmn.Move{}, err
	}

	p := g.board.GetPieceAt(from)

	m := cmn.Move{
		From:  from,
		To:    to,
		Piece: p,
	}

	toBB := cmn.Bitboard(1 << to)
	enemy := cmn.GetOpposite(g.board.SideToMove)
	if g.board.Occupancy[enemy]&toBB != 0 {
		m.Capture = g.board.GetPieceAt(to)
		m.Flags |= cmn.FlagCapture
	}

	if p == cmn.Pawn {
		diff := to - from

		if diff == 16 || diff == -16 {
			m.Flags |= cmn.FlagDoublePawnPush
		}

		if to == g.board.EnPassant {
			m.Flags |= cmn.FlagEnPassant
		}
	}

	if len(moveStr) == 5 {
		m.Promotion = cmn.PieceCharMap[rune(moveStr[4])].Type
		m.Flags |= cmn.FlagPromotion
	}

	if p == cmn.King {
		if (from == 4 && (to == 2 || to == 6)) ||
			(from == 60 && (to == 58 || to == 62)) {
			m.Flags |= cmn.FlagCastling
		}
	}

	return m, nil
}
