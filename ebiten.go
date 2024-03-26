package main

import (
	"syscall"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
}

func startEbiten() {
	// Set up ebiten
	ebiten.SetVsyncEnabled(true)
	ebiten.SetTPS(ebiten.SyncWithFPS)

	/* Set up our window */
	ebiten.SetWindowSize(defaultWindowWidth, defaultWindowHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("GoblinSiege")

	/* Start ebiten */
	if err := ebiten.RunGameWithOptions(newGame(), nil); err != nil {
		return
	}

	signalHandle <- syscall.SIGINT
}

func newGame() *Game {

	return &Game{}
}

/* Window size chaged, handle it */
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {

	return defaultWindowWidth, defaultWindowHeight
}
