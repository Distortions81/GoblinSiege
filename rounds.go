package main

import (
	"math/rand"
	"time"
)

func handleRounds() {
	for ServerRunning {
		votes.Lock.Lock()

		if votes.VoteState == VOTE_PLAYERS &&
			time.Since(votes.StartTime) > playerRoundTime {
			endVote()

		} else if votes.VoteState == VOTE_PLAYERS_DONE {
			votes.VoteState = VOTE_COMPUTER
			votes.CpuTime = time.Now()
			cpuTurn()

		} else if votes.VoteState == VOTE_COMPUTER &&
			time.Since(votes.CpuTime) > cpuRoundTime {

			votes.VoteState = VOTE_COMPUTER_DONE

		} else if votes.VoteState == VOTE_COMPUTER_DONE &&
			votes.GameRunning {
			if board.roundNum%3 == 0 {
				startVote()
			} else {
				votes.VoteState = VOTE_COMPUTER
				votes.CpuTime = time.Now()
				cpuTurn()
			}
		}

		if !votes.GameRunning {
			if time.Since(votes.RoundTime) > time.Second*5 {
				startGame()
			}
		}

		votes.Lock.Unlock()
		time.Sleep(time.Millisecond * 10)
	}
}

func cpuTurn() {
	board.lock.Lock()
	defer board.lock.Unlock()

	board.roundNum++

	if board.roundNum >= maxRounds {
		board.gameover = GAME_VICTORY
		endGame()
	}

	towerShootArrow()
	spawnGoblins()
	goblinAttack()

}

func spawnGoblins() {
	//Spawn goblins
	if board.roundNum%2 == 0 {
		goblin := getOtype("Goblin")
		rand := xyi{X: boardSizeX + enemyBoardX, Y: 1 + rand.Intn(boardSizeY-1)}
		if board.enemyMap[rand] == nil {
			board.enemyMap[rand] = &objectData{Pos: rand, oTypeP: goblin, Health: goblin.maxHealth, OldPos: xyi{X: rand.X + goblin.size.X, Y: rand.Y}}
		}
	}
}

func goblinAttack() {
	var newitems []*objectData
	//Detect defeat, defeat, do damage to towers, remove dead towers
	for _, item := range board.enemyMap {
		if item.dead {
			continue
		}

		//Detect game over
		oldItem := item

		//Setup next enemy position
		nextPos := item.Pos
		oldItem.OldPos = oldItem.Pos
		nextPos.X -= 1

		//Check towers and enemy positions before moving
		tower := board.playMap[nextPos]
		self := board.enemyMap[nextPos]
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
		delete(board.enemyMap, item.Pos)
		oldItem.Pos = nextPos
		newitems = append(newitems, oldItem)
	}

	//Add enemy back to new position
	for i, item := range newitems {
		board.enemyMap[item.Pos] = newitems[i]
		if item.Pos.X < 1 {
			board.gameover = GAME_DEFEAT
			endGame()
		}
	}
}
