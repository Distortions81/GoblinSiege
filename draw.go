package main

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func (g *Game) Draw(screen *ebiten.Image) {
	chatHistoryLock.Lock()
	buf := strings.Join(chatHistory, "")
	chatHistoryLock.Unlock()

	text.Draw(screen, buf, mplusNormalFont, 10, 25, color.White)
}
