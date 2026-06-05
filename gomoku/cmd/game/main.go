package main

import (
	"gomoku/internal/engine"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(1280, 900)
	ebiten.SetWindowTitle("Gomoku - Modern AI Board Game")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	gameEngine := engine.NewEngine()

	if err := ebiten.RunGame(gameEngine); err != nil {
		log.Fatal(err)
	}
}
