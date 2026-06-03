package game

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Engine struct {
	board *Board
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Loop() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		switch {
		case line == "uci":
			e.HandleUCI()

		case line == "isready":
			fmt.Println("readyok")

		case line == "ucinewgame":
			e.board, _ = NewBoard(StartPos)

		case strings.HasPrefix(line, "position"):
			e.HandlePosition(line)

		case strings.HasPrefix(line, "go"):
			e.HandleGo(line)

		case line == "quit":
			return
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
}

func (e *Engine) HandleUCI() {
	fmt.Println("id name Isabel")
	fmt.Println("id	author Elmago")
	fmt.Println("uciok")
}

func (e *Engine) HandlePosition(cmd string) {
	tokens := strings.Fields(cmd)

	e.board, _ = NewBoard(StartPos)

	idx := -1

	for i, t := range tokens {
		if t == "moves" {
			idx = i
			break
		}
	}

	if idx == -1 {
		return
	}

	for _, moveStr := range tokens[idx+1:] {
		move := e.parseMove(moveStr)
		e.board.MakeMove(move)
	}
}

func (e *Engine) HandleGo(cmd string) {
	depth := 5

	fields := strings.Fields(cmd)

	for i := range len(fields) - 1 {
		if fields[i] == "depth" {
			d, _ := strconv.Atoi(fields[i+1])
			depth = d
		}
	}

	move, err := e.board.FindBestMove(depth)

	if err != nil {
		fmt.Println("bestmove 0000")
		return
	}

	fmt.Printf("bestmove %s\n", move.String())
}

func (e *Engine) parseMove(m string) Move {
	legal := e.board.GenerateLegalMoves()

	for _, move := range legal {
		if move.String() == m {
			return move
		}
	}

	panic("illegal move")
}
