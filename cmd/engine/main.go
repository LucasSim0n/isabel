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
}
