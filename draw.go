package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func (g *Game) Draw(screen *ebiten.Image) {
	buf := fmt.Sprintf("FPS: %v", math.Round(ebiten.ActualFPS()))
	text.Draw(screen, buf, mplusNormalFont, 10, ScreenHeight-10, color.White)
}
