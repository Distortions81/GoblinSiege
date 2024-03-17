package main

import (
	"math/rand"
	"time"
)

func handleRounds() {
	for ServerRunning {
		UserMsgDict.Lock.Lock()

		if UserMsgDict.VoteState == VOTE_PLAYERS &&
			time.Since(UserMsgDict.StartTime) > playerRoundTime {
			endVote()

		} else if UserMsgDict.VoteState == VOTE_PLAYERS_DONE {
			UserMsgDict.VoteState = VOTE_COMPUTER
			UserMsgDict.StartTime = time.Now()
			cpuTurn()

		} else if UserMsgDict.VoteState == VOTE_COMPUTER &&
			time.Since(UserMsgDict.StartTime) > cpuRoundTime {

			UserMsgDict.VoteState = VOTE_COMPUTER_DONE

		} else if UserMsgDict.VoteState == VOTE_COMPUTER_DONE &&
			UserMsgDict.GameRunning {
			startVote()
		}

		UserMsgDict.Lock.Unlock()
		time.Sleep(time.Millisecond * 100)
	}
}

func cpuTurn() {
	board.lock.Lock()
	defer board.lock.Unlock()

	for _, item := range board.emap {
		oldItem := item
		if oldItem.Pos.X == 1 {
			board.gameover = GAME_DEFEAT
			endGame()
			return
		}
		nextPos := item.Pos
		nextPos.X -= 1

		tower := board.pmap[nextPos]
		self := board.emap[nextPos]
		if self != nil {
			continue
		}
		if tower != nil {
			if tower.Health >= 10 {
				tower.Health -= 10
			} else {
				delete(board.pmap, tower.Pos)
			}
			continue
		}
		delete(board.emap, item.Pos)
		oldItem.Pos = nextPos
		board.emap[oldItem.Pos] = oldItem
	}
	if board.roundNum == 0 || rand.Intn(3) == 0 {
		goblin := getOtype("Goblin")
		rand := xyi{X: boardSizeX + enemyBoardX, Y: 1 + rand.Intn(boardSizeY-1)}
		board.emap[rand] = &objectData{Pos: rand, oTypeP: goblin}
	}

	board.roundNum++
}
