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
			if board.roundNum%3 == 0 {
				startVote()
			} else {
				UserMsgDict.VoteState = VOTE_COMPUTER
				UserMsgDict.StartTime = time.Now()
				cpuTurn()
			}
		}

		UserMsgDict.Lock.Unlock()
		time.Sleep(time.Millisecond * 100)
	}
}

func cpuTurn() {
	board.lock.Lock()
	defer board.lock.Unlock()

	if board.roundNum >= maxRounds {
		board.gameover = GAME_VICTORY
		endGame()
	}

	towerShootArrow()
	goblinAttack()
	spawnGoblins()

	board.roundNum++
}

func spawnGoblins() {
	//Spawn goblins
	if board.roundNum%2 == 0 {
		goblin := getOtype("Goblin")
		rand := xyi{X: boardSizeX + enemyBoardX, Y: 1 + rand.Intn(boardSizeY-1)}
		if board.emap[rand] == nil {
			board.emap[rand] = &objectData{Pos: rand, oTypeP: goblin, Health: goblin.maxHealth}
		}
	}
}

func goblinAttack() {
	var newitems []*objectData
	//Detect defeat, defeat, do damage to towers, remove dead towers
	for _, item := range board.emap {
		if item.dead {
			continue
		}

		//Detect game over
		oldItem := item

		//Setup next enemy position
		nextPos := item.Pos
		nextPos.X -= 1

		//Check towers and enemy positions before moving
		tower := board.pmap[nextPos]
		self := board.emap[nextPos]
		if self != nil && !self.dead {
			continue
		}
		//If a tower is in our way, do damage
		if tower != nil && !tower.dead {
			tower.Health -= 10 + rand.Intn(10)
			if tower.Health <= 0 {
				tower.dead = true
			}
			continue
		}
		//Delete enemy, add to list
		delete(board.emap, item.Pos)
		oldItem.Pos = nextPos
		newitems = append(newitems, oldItem)
	}

	//Add enemy back to new position
	for i, item := range newitems {
		board.emap[item.Pos] = newitems[i]
		if item.Pos.X < 1 {
			board.gameover = GAME_DEFEAT
			endGame()
		}
	}
}
