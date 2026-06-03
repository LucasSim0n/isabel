package game

const TTSize = 1 << 20

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
	Move  Move
}

type TranspositionTable struct {
	Entries []TTEntry
}

func (tt *TranspositionTable) Store(hash uint64, depth int, score int, flag TTFlag, move Move) {
	index := hash & (TTSize - 1)

	tt.Entries[index] = TTEntry{
		Hash:  hash,
		Depth: depth,
		Score: score,
		Move:  move,
		Flag:  flag,
	}
}

func (tt *TranspositionTable) Probe(hash uint64, depth, alpha, beta int) (TTEntry, bool) {

	index := hash & (TTSize - 1)

	entry := tt.Entries[index]

	if entry.Hash != hash {
		return TTEntry{}, false
	}

	if entry.Depth < depth {
		return TTEntry{}, false
	}

	switch entry.Flag {
	case Exact:
		return entry, true

	case LowerBound:
		if entry.Score >= beta {
			return entry, true
		}

	case UpperBound:
		if entry.Score <= alpha {
			return entry, true
		}
	}
	return TTEntry{}, false
}

func (tt *TranspositionTable) Get(hash uint64) (TTEntry, bool) {
	i := hash & (TTSize - 1)
	e := tt.Entries[i]

	if e.Hash != hash {
		return TTEntry{}, false
	}

	return e, true
}
