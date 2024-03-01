package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func (g *Game) Draw(screen *ebiten.Image) {

	if UserMsgDict.Enabled {
		buf := fmt.Sprintf("Vote now! Votes: %v", UserMsgDict.Count)
		text.Draw(screen, buf, monoFont, 10, 30, color.White)
	}
}
