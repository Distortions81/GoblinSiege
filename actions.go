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
				go func() {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(200)))
					playVariated(SND_ARROW_SHOOT)
				}()

				if rand.Intn(2) != 0 {
					arrow := arrowData{tower: towerPos, target: enemy.Pos, missed: true, shot: curTime}
					board.arrowsShot = append(board.arrowsShot, arrow)
					break
				}
				arrow := arrowData{tower: towerPos, target: enemy.Pos, missed: false, shot: curTime}
				board.arrowsShot = append(board.arrowsShot, arrow)

				dmgAmt := 5 + rand.Intn(20)
				enemy.Health -= dmgAmt

				if enemy.Health <= 0 {
					board.enemyMap[enemy.Pos].dead = true
					board.enemyMap[enemy.Pos].diedAt = time.Now()
					go func() {
						time.Sleep(deathDelay + (time.Millisecond * time.Duration(rand.Intn(200))))
						playVariated(SND_GOBLIN_DIE)
					}()

					//For tweening
					board.enemyMap[enemy.Pos].OldPos = board.enemyMap[enemy.Pos].Pos
				}
				break
			}

		}
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
		item.attacking = false
		if tower != nil && !tower.dead {
			tower.Health -= 10 + rand.Intn(10)
			item.attacking = true

			go func() {
				time.Sleep(time.Millisecond*100 + time.Duration(rand.Intn(100)))
				playVariated(SND_AXE)
			}()
			if tower.Health <= 0 {
				go playSound(SND_TOWER_DIE)
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
