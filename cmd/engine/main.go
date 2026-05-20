package main

import (
	"chess/internal/board"
	"fmt"
	"log"
)

const (
	Kiwipete = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -"
	StartPos = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
)

func main() {
	b, err := board.NewBoard(Kiwipete)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(b.PerftDivide(5))
}
