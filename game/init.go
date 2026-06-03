package game

import "math/rand"

var knightAttacks [64]Bitboard
var kingAttacks [64]Bitboard

var pieceKeys [12][64]uint64
var castleKeys [16]uint64
var epKeys [8]uint64
var sideKey uint64

var TTable = TranspositionTable{
	Entries: make([]TTEntry, TTSize),
}

func generateKnightAttacks(sq int) Bitboard {
	var attacks Bitboard

	rank := sq / 8
	file := sq % 8

	moves := [8][2]int{
		{2, 1}, {1, 2}, {-1, 2}, {-2, 1},
		{-2, -1}, {-1, -2}, {1, -2}, {2, -1},
	}

	for _, m := range moves {
		r := rank + m[0]
		f := file + m[1]

		if r >= 0 && r < 8 && f >= 0 && f < 8 {
			attacks |= 1 << (r*8 + f)
		}
	}

	return attacks
}

func generateKingAttacks(sq int) Bitboard {
	var attacks Bitboard

	rank := sq / 8
	file := sq % 8

	moves := [8][2]int{
		{-1, 1}, {0, 1}, {1, 1},
		{-1, 0}, {1, 0},
		{-1, -1}, {0, -1}, {1, -1},
	}

	for _, m := range moves {
		r := rank + m[0]
		f := file + m[1]

		if r >= 0 && r < 8 && f >= 0 && f < 8 {
			attacks |= 1 << (r*8 + f)
		}
	}

	return attacks
}

func init() {
	for sq := range 64 {
		knightAttacks[sq] = generateKnightAttacks(sq)
		kingAttacks[sq] = generateKingAttacks(sq)
	}

	/*** Init Zobrist ***/
	rng := rand.New(rand.NewSource(1))

	for p := range 12 {
		for sq := range 64 {
			pieceKeys[p][sq] = rng.Uint64()
		}
	}

	for i := range 16 {
		castleKeys[i] = rng.Uint64()
	}

	for i := range 8 {
		epKeys[i] = rng.Uint64()
	}

	sideKey = rng.Uint64()
}
