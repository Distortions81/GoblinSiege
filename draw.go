package main

import (
	"goTwitchGame/sclean"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const pixelsPerLine = 36
const pad = 8
const namePixels = 170
const maxMsgLen = 300

const maxNameLen = 15
const chatLife = time.Second * 30

func (g *Game) Draw(screen *ebiten.Image) {
	chatHistoryLock.Lock()
	defer chatHistoryLock.Unlock()

	var justify int

	end := len(chatHistory)
	start := (end - maxShowLines)
	if start < 0 {
		start = 0
	}

	if end < maxShowLines {
		justify = (maxShowLines - end) * pixelsPerLine
	}

	z := 0
	for x := start; x < end; x++ {
		z++
		if time.Since(chatHistory[x].time) > chatLife {
			continue
		}
		name := sclean.TruncateString(chatHistory[x].sender, maxNameLen)
		namePix := text.BoundString(mplusNormalFont, name)

		maxAttempt := 1000
		count := 0
		lineLen := maxMsgLen
		for {
			count++

			msgPix := text.BoundString(mplusNormalFont, sclean.TruncateStringEllipsis(": "+chatHistory[x].message, lineLen))
			if msgPix.Max.X > (ScreenWidth - namePixels - pad) {
				lineLen--
			} else {
				break
			}

			if count > maxAttempt {
				lineLen = maxMsgLen
				break
			}
		}

		msg := sclean.TruncateStringEllipsis(": "+chatHistory[x].message, lineLen)

		text.Draw(screen, name, mplusNormalFont, (namePixels - namePix.Dx()), justify+(z*pixelsPerLine)-pad, chatHistory[x].color)
		text.Draw(screen, msg, mplusNormalFont, pad+namePixels, justify+(z*pixelsPerLine)-pad, color.White)
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
