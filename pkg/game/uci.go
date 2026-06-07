package game

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/LucasSim0n/isabel/pkg/board"
	"github.com/LucasSim0n/isabel/pkg/search"
)

func (g *GameManager) Loop() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		switch {
		case line == "uci":
			g.HandleUCI()

		case line == "isready":
			fmt.Println("readyok")

		case line == "ucinewgame":
			g.Board = board.NewStartPosGame()

		case strings.HasPrefix(line, "position"):
			g.HandlePosition(line)

		case strings.HasPrefix(line, "go"):
			g.HandleGo(line)

		case line == "quit":
			return
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
}

func (g *GameManager) HandleUCI() {
	fmt.Println("id name Isabel")
	fmt.Println("id	author Elmago")
	fmt.Println("uciok")
}

func (g *GameManager) HandlePosition(cmd string) {
	tokens := strings.Fields(cmd)

	idx := -1
	for i, t := range tokens {
		if t == "moves" {
			idx = i
			break
		}
	}

	switch tokens[1] {
	case "startpos":
		g.Board = board.NewStartPosGame()

	case "fen":
		var fen strings.Builder
		endFen := len(tokens)

		if idx != -1 {
			endFen = idx
		}
		for _, s := range tokens[2:endFen] {
			fmt.Fprintf(&fen, "%s ", s)
		}

		b, err := board.NewFromFen(fen.String())
		if err != nil {
			log.Fatalf("Error parsing fen: %s\n", err)
		}
		g.Board = b
	}

	if idx == -1 {
		return
	}

	for _, moveStr := range tokens[idx+1:] {
		move, _ := g.parseMove(moveStr)
		g.Board.MakeMove(move)
	}
}

func (g *GameManager) HandleGo(cmd string) {
	depth := 5

	fields := strings.Fields(cmd)

	for i := range len(fields) - 1 {
		if fields[i] == "depth" {
			d, _ := strconv.Atoi(fields[i+1])
			depth = d
			break
		}
	}

	move := g.Searcher.FindBestMove(g.Board, depth)
	fmt.Printf("bestmove %s\n", move.String())
	fmt.Println(search.Evaluate(g.Board))
}
