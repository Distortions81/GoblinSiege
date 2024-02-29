package main

import (
	"image/color"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const pixelsPerLine = 36

func (g *Game) Draw(screen *ebiten.Image) {
	chatHistoryLock.Lock()
	buf := strings.Join(chatHistory, "")
	chatHistoryLock.Unlock()

	text.Draw(screen, buf, mplusNormalFont, 10, pixelsPerLine, color.White)
}

func adjMaxLines() {
	chatHistoryLock.Lock()
	maxLines = ScreenHeight / pixelsPerLine
	log.Printf("Max lines set to: %v", maxLines)
	chatHistoryLock.Unlock()
}
