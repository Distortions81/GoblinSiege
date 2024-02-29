package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const pixelsPerLine = 36
const vpad = 8

func (g *Game) Draw(screen *ebiten.Image) {
	chatHistoryLock.Lock()
	defer chatHistoryLock.Unlock()

	end := len(chatHistory)
	start := (end - maxShowLines)
	if start < 0 {
		start = 0
	}

	z := 0
	for x := start; x < end; x++ {
		z++
		text.Draw(screen, chatHistory[x].message, mplusNormalFont, 10, (z*pixelsPerLine)-vpad, chatHistory[x].color)
	}
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
