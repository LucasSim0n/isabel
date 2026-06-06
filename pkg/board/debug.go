package board

import "fmt"

func (b *Board) Perft(depth int) uint64 {
	if depth == 0 {
		return 1
	}

	var nodes uint64

	moves := b.GenerateLegalMoves()

	for _, move := range moves {

		before := b.Clone()

		undo := b.MakeMove(move)

		nodes += b.Perft(depth - 1)

		b.UnmakeMove(move, undo)

		if !before.Equals(b) {
			fmt.Printf("Corruption after: %s depth=%d\n", move, depth)

			fmt.Println("BEFORE:")
			fmt.Println(before)

			fmt.Println("AFTER:")
			fmt.Println(b)
		}
	}

	return nodes
}

func (b *Board) PerftDivide(depth int) uint64 {
	if depth == 0 {
		return 1
	}

	var total uint64

	moves := b.GenerateLegalMoves()

	for _, move := range moves {
		undo := b.MakeMove(move)

		nodes := b.Perft(depth - 1)

		b.UnmakeMove(move, undo)

		fmt.Printf("%s: %d\n", move.String(), nodes)

		total += nodes
	}

	return total
}

func (b *Board) Equals(other *Board) bool {

	for i := range 12 {
		if b.Pieces[i] != other.Pieces[i] {
			fmt.Printf("Pieces[%d] mismatch\n", i)
			return false
		}
	}

	for i := range 3 {
		if b.Occupancy[i] != other.Occupancy[i] {
			fmt.Printf("Occupancy[%d] mismatch\n", i)
			return false
		}
	}

	if b.SideToMove != other.SideToMove {
		fmt.Println("SideToMove mismatch")
		return false
	}

	if b.CastlingRights != other.CastlingRights {
		fmt.Println("CastlingRights mismatch")
		return false
	}

	if b.EnPassant != other.EnPassant {
		fmt.Println("EnPassant mismatch")
		return false
	}

	if b.HalfmoveClock != other.HalfmoveClock {
		fmt.Println("HalfmoveClock mismatch")
		return false
	}

	if b.FullmoveNumber != other.FullmoveNumber {
		fmt.Println("FullmoveNumber mismatch")
		return false
	}

	if b.KingSq != other.KingSq {
		fmt.Println("KingSq mismatch")
		return false
	}

	if b.Hash != other.Hash {
		fmt.Println("Hash mismatch")
		return false
	}

	return true
}
