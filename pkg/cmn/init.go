package cmn

import (
	"math/rand"
)

func init() {
	for sq := range 64 {
		KnightAttacks[sq] = generateKnightAttacks(sq)
		KingAttacks[sq] = generateKingAttacks(sq)
	}

	/*** Init Zobrist ***/
	rng := rand.New(rand.NewSource(1))

	for p := range 12 {
		for sq := range 64 {
			PieceKeys[p][sq] = rng.Uint64()
		}
	}

	for i := range 16 {
		CastleKeys[i] = rng.Uint64()
	}

	for i := range 8 {
		EpKeys[i] = rng.Uint64()
	}

	SideKey = rng.Uint64()
}
