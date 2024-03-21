package main

import (
	"log"
	"math/rand"
	"time"
)

func addTower() {

	if votes.VoteCount > 0 &&
		votes.Result.X > 0 &&
		votes.Result.Y > 0 &&
		votes.Result.X <= boardSizeX &&
		votes.Result.Y <= boardSizeY {

		tpos := votes.Result
		if board.enemyMap[tpos] == nil && board.playMap[tpos] == nil {
			board.playMap[tpos] = &objectData{Pos: tpos, sheetP: &obj_tower1, Health: obj_tower1.health}
		} else {
			log.Println("COLLISION!")
		}
	} else {

		log.Println("Not enough votes, picking random.")
		//Invalid or not enough votes, pick a pos at random
		tpos := xyi{X: rand.Intn(boardSizeX-1) + 1, Y: rand.Intn(boardSizeY-1) + 1}
		if board.enemyMap[tpos] == nil && board.playMap[tpos] == nil {
			board.playMap[tpos] = &objectData{Pos: tpos, sheetP: &obj_tower1, Health: obj_tower1.health, aniOffset: uint64(rand.Intn(obj_tower1.frames))}
		}
	}

}

func towerShootArrow() {
	curTime := time.Now()

	var didShoot, didDie, didHurt int

	for _, item := range board.playMap {
		if item.dead {
			continue
		}

		//Shoot from tower top
		towerPos := item.Pos
		towerPos.Y -= 1

		for _, enemy := range board.enemyMap {
			if enemy.dead {
				continue
			}

			//If enemy within range
			if Distance(item.Pos, enemy.Pos) < 6 {
				didShoot++

				if rand.Intn(2) != 0 {
					arrow := arrowData{tower: towerPos, target: enemy.Pos, missed: true, shot: curTime}
					board.arrowsShot = append(board.arrowsShot, arrow)
					break
				}
				arrow := arrowData{tower: towerPos, target: enemy.Pos, missed: false, shot: curTime}
				board.arrowsShot = append(board.arrowsShot, arrow)

				dmgAmt := 5 + rand.Intn(20)
				enemy.Health -= dmgAmt
				didHurt++
				if enemy.Health <= 0 {
					board.enemyMap[enemy.Pos].dead = true
					board.enemyMap[enemy.Pos].diedAt = time.Now()
					didDie++

					//For tweening
					board.enemyMap[enemy.Pos].OldPos = board.enemyMap[enemy.Pos].Pos
				}
				break
			}

		}
	}

	if didShoot > 0 {
		go playVariated(SND_ARROW_SHOOT, didShoot)
	}

	if didDie > 0 {
		go func() {
			time.Sleep(deathDelay)
			playSound(SND_GOBLIN_DIE)
		}()
	}

	if didHurt > 0 {
		//
	}
}

func spawnGoblins() {
	//Spawn goblins
	if board.moveNum%2 == 0 {
		rpos := xyi{X: boardSizeX + enemyBoardX, Y: 1 + rand.Intn(boardSizeY-1)}
		if board.enemyMap[rpos] == nil {
			board.enemyMap[rpos] = &objectData{Pos: rpos, sheetP: &obj_goblinBarb, Health: obj_goblinBarb.health, OldPos: xyi{X: rpos.X, Y: rpos.Y}, aniOffset: uint64(rand.Intn(obj_goblinBarb.frames))}
		}
	}
}

func goblinAttack() {
	var didHit int

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
			didHit++
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

	if didHit > 0 {
		playVariated(SND_AXE, didHit)
	}
}
