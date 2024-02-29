package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func (g *Game) Draw(screen *ebiten.Image) {

	if votesEnabled {
		buf := fmt.Sprintf("Vote now! Votes: %v", commandCount)
		text.Draw(screen, buf, monoFont, 10, 30, color.White)
	}
}
