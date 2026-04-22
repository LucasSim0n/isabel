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

	fmt.Println(b.Perft(1))
	fmt.Println(b.Perft(2))
	fmt.Println(b.Perft(3))
	fmt.Println(b.Perft(4))
	fmt.Println(b.Perft(5))
}
