package main

import (
	"chess/internal/board"
	"fmt"
)

func main() {
	b, err := board.NewBoard("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(b.String())

	moves := b.GenerateMoves()

	for _, m := range *moves {
		if m.Piece == board.Knight {
			fmt.Printf("%d\n", m.To)
		}
	}
}
