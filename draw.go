package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func (g *Game) Draw(screen *ebiten.Image) {

	UserMsgDict.Lock.Lock()
	defer UserMsgDict.Lock.Unlock()

	if UserMsgDict.Enabled {
		buf := fmt.Sprintf("Vote now! Votes: %v", UserMsgDict.Count)
		text.Draw(screen, buf, monoFont, 10, 30, color.White)
	} else if UserMsgDict.Count > 0 {
		buf := fmt.Sprintf("Result: %v,%v", UserMsgDict.Result.X, UserMsgDict.Result.Y)
		text.Draw(screen, buf, monoFont, 10, 30, color.White)
	} else {
		text.Draw(screen, "Not enough votes.", monoFont, 10, 30, color.White)
	}
}
