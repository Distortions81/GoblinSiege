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

	end := len(chatHistory)
	start := (end - maxShowLines)
	if start < 0 {
		start = 0
	}

	showLines := chatHistory[start:end]
	buf := strings.Join(showLines, "")

	vertPad := pixelsPerLine
	if numLines < maxShowLines {
		vertPad = ((maxShowLines - numLines) * pixelsPerLine)
	}

	chatHistoryLock.Unlock()

	text.Draw(screen, buf, mplusNormalFont, 10, vertPad, color.White)
}

func adjMaxLines() {
	chatHistoryLock.Lock()

	oldMaxLines := maxShowLines
	maxShowLines = (ScreenHeight / pixelsPerLine)
	if maxShowLines != oldMaxLines {
		log.Printf("Max lines set to: %v", maxShowLines)
		trimChatHistory()
	}
	chatHistoryLock.Unlock()
}
