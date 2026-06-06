package board

import "github.com/LucasSim0n/isabel/pkg/cmn"

var startPosGame *Board

func init() {
	startPosGame, _ = NewFromFen(cmn.StartPos)
}
