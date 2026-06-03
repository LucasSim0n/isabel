package main

import "chess/game"

func main() {
	e := game.NewEngine()
	e.Loop()
}
