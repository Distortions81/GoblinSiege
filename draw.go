package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hako/durafmt"
)

func (g *Game) Draw(screen *ebiten.Image) {

	UserMsgDict.Lock.Lock()
	defer UserMsgDict.Lock.Unlock()

	if UserMsgDict.Enabled && time.Since(UserMsgDict.StartTime) > roundTime {
		endVote()
	}

	if UserMsgDict.Enabled {
		buf := fmt.Sprintf("Vote now: !letter number\nVotes: %v (%v left)", UserMsgDict.Count, durafmt.Parse(time.Until(UserMsgDict.StartTime.Add(roundTime))).LimitFirstN(1))
		text.Draw(screen, buf, monoFont, 10, 30, color.White)
		return
	} else {
		if UserMsgDict.Count > 0 {
			buf := fmt.Sprintf("Result: %v,%v", UserMsgDict.Result.X, UserMsgDict.Result.Y)
			text.Draw(screen, buf, monoFont, 10, 30, color.White)
		} else if time.Since(UserMsgDict.StartTime) < time.Second*10 {
			text.Draw(screen, "Not enough votes, stopped.", monoFont, 10, 30, color.White)
		}
	}

	drawGameBoard(screen)
}
