package tt

import "github.com/LucasSim0n/isabel/pkg/cmn"

type TTFlag uint8

const (
	Exact TTFlag = iota
	LowerBound
	UpperBound
)

type TTEntry struct {
	Hash  uint64
	Flag  TTFlag
	Depth int
	Score int
	Move  cmn.Move
}
