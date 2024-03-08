package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hako/durafmt"
)

func (g *Game) Draw(screen *ebiten.Image) {

	drawGameBoard(screen)

	UserMsgDict.Lock.Lock()
	defer UserMsgDict.Lock.Unlock()

	if UserMsgDict.Voting &&
		time.Since(UserMsgDict.StartTime) > roundTime {
		endVote()
	}
	if UserMsgDict.GameRunning &&
		time.Since(UserMsgDict.StartTime) > restTime+roundTime {
		startVote()
	}

	if UserMsgDict.Voting {
		vector.DrawFilledRect(screen, 0, 0, 400, 100, ColorSmoke, true)
		buf := fmt.Sprintf("Vote now: %v1 1 -- (x,y 1-25)\nVotes: %v\n%v remaining...", userSettings.CmdPrefix, UserMsgDict.Count, durafmt.Parse(time.Until(UserMsgDict.StartTime.Add(roundTime)).Round(time.Second)).LimitFirstN(1))
		text.Draw(screen, buf, monoFont, 10, 30, color.White)
	} else {
		if UserMsgDict.Count > 0 {
			buf := fmt.Sprintf("Result: %v,%v", UserMsgDict.Result.X, UserMsgDict.Result.Y)
			text.Draw(screen, buf, monoFont, 10, 30, color.White)
		} else if !UserMsgDict.GameRunning {
			text.Draw(screen, "No game active.", monoFont, 10, 30, color.White)
		}

	}
}
