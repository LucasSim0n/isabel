package board

var knightAttacks [64]Bitboard
var kingAttacks [64]Bitboard

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
		f := file + m[0]

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
}
