package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func startEbiten() {
	// Set up ebiten
	ebiten.SetVsyncEnabled(true)
	ebiten.SetTPS(ebiten.SyncWithFPS)

	/* We manaually clear, so we aren't forced to draw every frame */
	ScreenWidth, ScreenHeight = defaultWindowWidth, defaultWindowHeight

	/* Set up our window */
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("goTwitchGame")

	/* Start ebiten */
	if err := ebiten.RunGameWithOptions(newGame(), nil); err != nil {
		return
	}
}

func newGame() *Game {

	return &Game{}
}

/* Window size chaged, handle it */
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {

	if outsideWidth != ScreenWidth || outsideHeight != ScreenHeight {
		ScreenWidth, ScreenHeight = outsideWidth, outsideHeight
	}

	return ScreenWidth, ScreenHeight
}

func (g *Game) Update() error {
	return nil
}