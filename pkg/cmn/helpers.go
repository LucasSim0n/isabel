package cmn

import "fmt"

func NotationToSquare(notation string) (int, error) {
	if len(notation) != 2 {
		return 0, fmt.Errorf("invalid square notation: %s", notation)
	}

	fileChar := notation[0]
	rankChar := notation[1]

	file := int(fileChar - 'a')
	rank := int(rankChar - '1')

	if file < 0 || file > 7 || rank < 0 || rank > 7 {
		return 0, fmt.Errorf("square notation out of bounds: %s", notation)
	}

	return rank*8 + file, nil
}

func SquareToNotation(sq int) string {
	if sq < 0 || sq > 63 {
		return "-"
	}
	file := rune('a' + (sq % 8))
	rank := rune('1' + (sq / 8))
	return string([]rune{file, rank})
}

func GetOpposite(c Color) Color {
	if c == White {
		return Black
	}
	return White
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
